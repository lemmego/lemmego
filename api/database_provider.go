package api

import (
	"context"

	"github.com/lemmego/lemmego/api/db"
)

type DatabaseServiceProvider struct {
	*BaseServiceProvider
}

func (provider *DatabaseServiceProvider) Register(app *App) {
	dbConfig := &db.Config{
		Driver:   Config("db.driver").(string),
		Host:     Config("db.host").(string),
		Port:     Config("db.port").(int),
		Database: Config("db.database").(string),
		User:     Config("db.username").(string),
		Password: Config("db.password").(string),
	}

	dbc, err := db.NewConnection(dbConfig).
		// WithForceCreateDb(). // Force create db if not exists
		Open()
	if err != nil {
		panic(err)
	}
	//app.Bind((*db.DB)(nil), func() *db.DB {
	//	return dbc
	//})
	app.db = dbc
	app.dbFunc = func(c context.Context, config *db.Config) (*db.DB, error) {
		if config == nil {
			config = dbConfig
		}
		return db.NewConnection(config).
			// WithForceCreateDb(). // Force create db if not exists
			Open()
	}
}

func (provider *DatabaseServiceProvider) Boot() {
	//
}
