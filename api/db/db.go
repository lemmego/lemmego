package db

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"

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
)

type DB struct {
	*gorm.DB
}

type Model struct {
	gorm.Model
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

	return &DB{db}, nil
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

	return &DB{db}, nil
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
		return c.connectToMySQL()
	case "postgres":
		return c.connectToPostgres()
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
