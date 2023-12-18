package fluent

import (
	"errors"
	"log"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
)

type DBSession interface {
	db.Session
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func ConnectDB(dialect string, dbConfig *DBConfig) (DBSession, error) {
	if dialect == "postgres" {
		settings := postgresql.ConnectionURL{
			Database: dbConfig.Database,
			Host:     dbConfig.Host,
			User:     dbConfig.User,
			Password: dbConfig.Password,
		}
		sess, err := postgresql.Open(settings)
		if err != nil {
			log.Fatal("postgresql.Open: ", err)
		}
		return sess, err
	}

	if dialect == "mysql" {
		settings := mysql.ConnectionURL{
			Database: dbConfig.Database,
			Host:     dbConfig.Host,
			User:     dbConfig.User,
			Password: dbConfig.Password,
		}
		sess, err := mysql.Open(settings)
		if err != nil {
			log.Fatal("mysql.Open: ", err)
		}
		return sess, err
	}
	return nil, errors.New("Invalid dialect")
}
