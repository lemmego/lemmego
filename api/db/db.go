package db

import (
	"errors"
	"fmt"
	"log"
	"pressebo/api/logger"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
)

var (
	ErrUnknownDriver  = errors.New("unknown driver")
	ErrNoSuchDatabase = errors.New("no such database exists")
)

type DBSession interface {
	db.Session
}

type Cond = db.Cond

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Connection struct {
	forceCreateDb bool
	config        *DBConfig
	session       DBSession
}

func NewConnection(dbc *DBConfig) *Connection {
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
	if c.session == nil {
		return false
	}

	if err := c.session.Ping(); err != nil {
		return false
	}

	return true
}

func (c *Connection) WithDatabase(database string) *Connection {
	c.config.Database = database
	return c
}

func (c *Connection) connectToMySQL() (DBSession, error) {
	dbConfig := c.config
	settings := mysql.ConnectionURL{
		Host:     dbConfig.Host,
		User:     dbConfig.User,
		Password: dbConfig.Password,
		Database: dbConfig.Database,
	}
	sess, err := mysql.Open(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	logger.Log().Info(fmt.Sprintf("created db session %s", sess.Name()))
	return sess, nil
}

func (c *Connection) connectToPostgres() (DBSession, error) {
	dbConfig := c.config
	settings := postgresql.ConnectionURL{
		Host:     dbConfig.Host,
		User:     dbConfig.User,
		Password: dbConfig.Password,
		Database: dbConfig.Database,
	}

	sess, err := postgresql.Open(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	logger.Log().Info(fmt.Sprintf("created db session %s", sess.Name()))
	return sess, nil
}

func (c *Connection) Close() error {
	return c.session.Close()
}

func (c *Connection) existsDb() error {
	var sess db.Session
	var err error
	dbConfig := c.config
	database := c.config.Database

	defer func() {
		c.WithDatabase(database)
		if sess != nil {
			logger.Log().Info(fmt.Sprintf("closing db session %s", sess.Name()))
			sess.Close()
		}
	}()

	if dbConfig.Driver == "postgres" {
		sess, err = c.WithDatabase("postgres").connectToPostgres()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		if row, err := sess.SQL().QueryRow("SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('" + database + "')"); err != nil {
			return fmt.Errorf("failed to check db existence: %w", err)
		} else {
			var fetchedDatabase string
			row.Scan(&fetchedDatabase)
			if fetchedDatabase != database {
				return ErrNoSuchDatabase
			}
			return nil
		}
	}

	if dbConfig.Driver == "mysql" {
		sess, err = c.WithDatabase("mysql").connectToMySQL()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		if row, err := sess.SQL().QueryRow("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '" + database + "'"); err != nil {
			return fmt.Errorf("failed to check db existence: %w", err)
		} else {
			var fetchedDatabase string
			row.Scan(&fetchedDatabase)
			if fetchedDatabase != database {
				return ErrNoSuchDatabase
			}
			return nil
		}
	}

	return ErrUnknownDriver
}

func (c *Connection) Open() (DBSession, error) {
	if c.IsOpen() {
		log.Println(fmt.Sprintf("closing db session %s, before opening a new one", c.session.Name()))
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
	var sess db.Session
	var err error

	defer func() {
		c.WithDatabase(database)
		if sess != nil {
			logger.Log().Info(fmt.Sprintf("closing db session %s", sess.Name()))
			sess.Close()
		}
	}()

	if dbConfig.Driver == "postgres" {
		if sess, err = c.WithDatabase("postgres").Open(); err != nil {
			return err
		}
		res, err := sess.SQL().Exec("CREATE DATABASE " + database + " WITH OWNER " + dbConfig.User)
		if err != nil {
			return err
		}
		if res != nil {
			log.Println("database", database, "created")
			return nil
		}
	}

	if dbConfig.Driver == "mysql" {
		if sess, err = c.WithDatabase("mysql").Open(); err != nil {
			return err
		}

		res, err := sess.SQL().Exec("CREATE DATABASE IF NOT EXISTS " + database)
		if err != nil {
			return err
		}
		if res != nil {
			log.Println("database", database, "created")
			return nil
		}
	}

	return nil
}
