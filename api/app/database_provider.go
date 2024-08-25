package app

import (
	"context"

	"github.com/lemmego/lemmego/api/config"
	"github.com/lemmego/lemmego/api/db"
)

type DatabaseServiceProvider struct {
	*BaseServiceProvider
}

func (provider *DatabaseServiceProvider) Register(app *App) {
	dbConfig := &db.Config{
		ConnName: "default",
		Driver:   config.Get[string]("database.connections.default.driver"),
		Host:     config.Get[string]("database.connections.default.host"),
		Port:     config.Get[int]("database.connections.default.port"),
		Database: config.Get[string]("database.connections.default.database"),
		User:     config.Get[string]("database.connections.default.user"),
		Password: config.Get[string]("database.connections.default.password"),
		Params:   config.Get[string]("database.connections.default.params"),
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
	app.SetDB(dbc)
	app.SetDbFunc(func(c context.Context, config *db.Config) (*db.DB, error) {
		if config == nil {
			config = dbConfig
		}
		return db.NewConnection(config).
			// WithForceCreateDb(). // Force create db if not exists
			Open()
	})
}

func (provider *DatabaseServiceProvider) Boot() {
	//
}
