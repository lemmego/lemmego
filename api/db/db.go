package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DefaultPostgresDB = "postgres"
	DefaultMysqlDB    = "mysql"
)

var (
	ErrUnknownDriver            = errors.New("unknown driver")
	ErrNoSuchDatabase           = errors.New("no such database exists")
	ErrCannotConnectToDefaultDB = errors.New("cannot connect to default db")

	dbInstances = make(map[string]*DB)
)

type Resolver struct {
	Ctx context.Context
	New func(context.Context) (*DB, error)
}

type ResolverFunc func(context.Context) (*DB, error)

type DB struct {
	*gorm.DB
}

func (db *DB) BindWhere(c context.Context, columnName string) *gorm.DB {
	return db.Where(fmt.Sprintf("%s = ?", columnName), c.Value(columnName))
}

type Model struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func Resolve(ctx context.Context, resolver ResolverFunc) (*DB, error) {
	db, err := resolver(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Determine if the map has an entry with the key
func Has(connectionId string) bool {
	_, ok := dbInstances[connectionId]
	return ok
}

func Get(arg ...interface{}) *DB {
	if len(arg) == 0 {
		if val, ok := dbInstances["default"]; ok {
			return val
		}
		return nil
	}

	// Check if arg[0] is of type string
	if val, ok := arg[0].(string); ok {
		if val == "default" {
			if val, ok := dbInstances["default"]; ok {
				return val
			}
			return nil
		}
		if val, ok := dbInstances[val]; ok {
			return val
		}
		return nil
	}

	return nil
}

func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()

	if err != nil {
		return errors.New("failed to close db connection")
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}
	return nil
}

type Config struct {
	ConnName string
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Params   string
}

type Connection struct {
	forceCreateDb bool
	config        *Config
	db            *DB
}

func NewConnection(dbc *Config) *Connection {
	if dbc.ConnName == "" {
		dbc.ConnName = "default"
	}

	if dbc.Driver != "mysql" && dbc.Driver != "postgres" && dbc.Driver != "sqlite" {
		panic("unsupported driver")
	}

	if dbc.Driver != "sqlite" && dbc.Host == "" {
		dbc.Host = "localhost"
	}

	if dbc.Driver == "mysql" && dbc.Port == 0 {
		dbc.Port = 3306
	}

	if dbc.Driver == "postgres" && dbc.Port == 0 {
		dbc.Port = 5432
	}

	if dbc.Driver != "sqlite" && dbc.User == "" {
		panic("db username must be provided")
	}

	if dbc.Driver == "sqlite" && dbc.Database == "" {
		panic("a path to database file must be provided for sqlite")
	}

	return &Connection{false, dbc, nil}
}

func (c *Connection) WithForceCreateDb() *Connection {
	c.forceCreateDb = true
	return c
}

func (c *Connection) IsOpen() bool {
	if c.db == nil {
		return false
	}

	sqlDB, err := c.db.DB.DB()

	if err != nil {
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		return false
	}

	return true
}

func (c *Connection) WithDatabase(database string) *Connection {
	c.config.Database = database
	return c
}

func (c *Connection) connectToMySQL() (*DB, error) {
	dbConfig := c.config
	dsn := &DataSource{
		Dialect:  DialectMySQL,
		Host:     dbConfig.Host,
		Port:     strconv.Itoa(dbConfig.Port),
		Username: dbConfig.User,
		Password: dbConfig.Password,
		Name:     dbConfig.Database,
		Params:   dbConfig.Params,
	}
	dsnStr, err := dsn.String()
	db, err := gorm.Open(mysql.Open(dsnStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	slog.Info(fmt.Sprintf("created db session %s", dsn.Name))

	dbInstance := &DB{db}
	c.db = dbInstance
	return dbInstance, nil
}

func (c *Connection) connectToPostgres() (*DB, error) {
	dbConfig := c.config
	dsn := &DataSource{
		Dialect:  DialectPostgres,
		Host:     dbConfig.Host,
		Port:     strconv.Itoa(dbConfig.Port),
		Username: dbConfig.User,
		Password: dbConfig.Password,
		Name:     dbConfig.Database,
		Params:   dbConfig.Params,
	}
	dsnStr, err := dsn.String()
	db, err := gorm.Open(postgres.Open(dsnStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	slog.Info(fmt.Sprintf("created db session %s", dsn.Name))

	dbInstance := &DB{db}
	c.db = dbInstance
	return dbInstance, nil
}

func (c *Connection) Close() error {
	return c.db.Close()
}

func (c *Connection) existsDb() error {
	var dbConn *DB
	var err error
	dbConfig := c.config
	database := c.config.Database

	defer func() {
		c.WithDatabase(database)
		if dbConn != nil {
			slog.Info(fmt.Sprintf("closing db session %s", dbConn.Name()))
			dbConn.Close()
		}
	}()

	if dbConfig.Driver == DialectPostgres {
		dbConn, err = c.WithDatabase(DefaultPostgresDB).connectToPostgres()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		var fetchedDatabase string
		dbConn.Raw("SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower(?)", database).Scan(&fetchedDatabase)
		if fetchedDatabase == "" {
			return ErrCannotConnectToDefaultDB
		}
		if fetchedDatabase != database {
			return ErrNoSuchDatabase
		}
		return nil
	}

	if dbConfig.Driver == DialectMySQL {
		dbConn, err = c.WithDatabase(DefaultMysqlDB).connectToMySQL()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		var fetchedDatabase string
		dbConn.Raw("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", database).Scan(&fetchedDatabase)
		if fetchedDatabase == "" {
			return ErrCannotConnectToDefaultDB
		}
		if fetchedDatabase != database {
			return ErrNoSuchDatabase
		}
		return nil
	}

	return ErrUnknownDriver
}

func (c *Connection) Open() (*DB, error) {
	if c.IsOpen() {
		log.Println(fmt.Sprintf("closing db session %s, before opening a new one", c.config.Database))
		c.Close()
	}

	if c.forceCreateDb {
		err := c.existsDb()
		if err != nil && errors.Is(err, ErrNoSuchDatabase) {
			c.createDb()
		}
	}

	switch c.config.Driver {
	case "mysql":
		db, err := c.connectToMySQL()
		if err != nil {
			return nil, err
		}
		dbInstances[c.config.ConnName] = db
		return db, nil
	case "postgres":
		db, err := c.connectToPostgres()
		if err != nil {
			return nil, err
		}
		dbInstances[c.config.ConnName] = db
		return db, nil
	default:
		return nil, ErrUnknownDriver
	}
}

func (c *Connection) createDb() error {
	dbConfig := c.config
	database := dbConfig.Database
	var db *DB
	var err error

	defer func() {
		c.WithDatabase(database)
		if db != nil {
			slog.Info(fmt.Sprintf("closing db session %s", database))
			db.Close()
		}
	}()

	if dbConfig.Driver == "postgres" {
		if db, err = c.WithDatabase(DefaultPostgresDB).Open(); err != nil {
			return err
		}
		err := db.Exec("CREATE DATABASE " + database + " WITH OWNER " + dbConfig.User).Error
		if err != nil {
			return err
		} else {
			slog.Info("database", database, "created")
			return nil
		}
	}

	if dbConfig.Driver == "mysql" {
		if db, err = c.WithDatabase("mysql").Open(); err != nil {
			return err
		}

		err := db.Exec("CREATE DATABASE IF NOT EXISTS " + database).Error
		if err != nil {
			return err
		} else {
			log.Println("database", database, "created")
			return nil
		}
	}

	return nil
}
