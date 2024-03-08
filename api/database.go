package api

import (
	"errors"
	"fmt"
	"log"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
)

var (
	ErrUnknownDialect = errors.New("unknown dialect")
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

func DBExists(dialect string, dbConfig *DBConfig) (bool, error) {
	var sess db.Session
	var err error

	defer func() {
		if sess != nil {
			log.Println("closing db session", sess.Name())
			sess.Close()
		}
	}()

	if dialect == "postgres" {
		settings := postgresql.ConnectionURL{
			Host:     dbConfig.Host,
			User:     dbConfig.User,
			Password: dbConfig.Password,
		}

		sess, err = postgresql.Open(settings)
		if err != nil {
			return false, fmt.Errorf("failed to connect: %w", err)
		}

		if row, err := sess.SQL().QueryRow("SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('" + dbConfig.Database + "')"); err != nil {
			return false, fmt.Errorf("failed to check db existence: %w", err)
		} else {
			var database string
			row.Scan(&database)
			if database != dbConfig.Database {
				return false, ErrNoSuchDatabase
			}
			return true, nil
		}
	}

	if dialect == "mysql" {
		settings := mysql.ConnectionURL{
			Host:     dbConfig.Host,
			User:     dbConfig.User,
			Password: dbConfig.Password,
			Database: "mysql",
		}

		sess, err = mysql.Open(settings)
		if err != nil {
			fmt.Println(settings.Host, settings.User, settings.Password, settings.Database)
			return false, fmt.Errorf("failed to connect: %w", err)
		}

		if row, err := sess.SQL().QueryRow("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '" + dbConfig.Database + "'"); err != nil {
			return false, fmt.Errorf("failed to check db existence: %w", err)
		} else {
			var database string
			row.Scan(&database)
			if database != dbConfig.Database {
				return false, ErrNoSuchDatabase
			}
			return true, nil
		}
	}

	return false, ErrUnknownDialect
}

func CreateDB(dialect string, dbConfig *DBConfig) error {
	var sess db.Session
	var err error

	defer func() {
		if sess != nil {
			log.Println("closing db session", sess.Name())
			sess.Close()
		}
	}()

	if sess, err = ConnectToDefaultDb(dialect, dbConfig); err != nil {
		return err
	}

	if dialect == "postgres" {
		res, err := sess.SQL().Exec("CREATE DATABASE " + dbConfig.Database + " WITH OWNER " + dbConfig.User)
		if err != nil {
			return err
		}
		if res != nil {
			return nil
		}
	}

	if dialect == "mysql" {
		res, err := sess.SQL().Exec("CREATE DATABASE IF NOT EXISTS " + dbConfig.Database)
		if err != nil {
			return err
		}
		if res != nil {
			log.Println("database", dbConfig.Database, "created")
			return nil
		}
	}

	return nil
}

func ConnectToDefaultDb(dialect string, dbConfig *DBConfig) (DBSession, error) {
	if dialect == "postgres" {
		settings := postgresql.ConnectionURL{
			Host:     dbConfig.Host,
			User:     dbConfig.User,
			Password: dbConfig.Password,
		}

		sess, err := postgresql.Open(settings)
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}

		log.Println("created db session", sess.Name())
		return sess, nil
	}

	if dialect == "mysql" {
		settings := mysql.ConnectionURL{
			Host:     dbConfig.Host,
			User:     dbConfig.User,
			Password: dbConfig.Password,
			Database: "mysql",
		}

		sess, err := mysql.Open(settings)
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}

		log.Println("created db session", sess.Name())
		return sess, nil
	}

	return nil, ErrUnknownDialect
}

func ConnectToDatabase(dialect string, dbConfig *DBConfig) (DBSession, error) {
	if dialect == "postgres" {
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

		log.Println("created db session", sess.Name())
		return sess, nil
	}

	if dialect == "mysql" {
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

		log.Println("created db session", sess.Name())
		return sess, nil
	}

	return nil, ErrUnknownDialect
}

func ConnectDB(dialect string, dbConfig *DBConfig) (DBSession, error) {
	if exists, _ := DBExists(dialect, dbConfig); !exists {
		log.Println("database", dbConfig.Database, "does not exist, creating it")
		if err := CreateDB(dialect, dbConfig); err != nil {
			return nil, err
		}
	}

	if sess, err := ConnectToDatabase(dialect, dbConfig); err != nil {
		return nil, err
	} else {
		return sess, nil
	}
}
