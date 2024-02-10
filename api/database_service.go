package api

import (
	"os"
)

type DatabaseServiceProvider struct {
	BaseServiceProvider
}

func (provider *DatabaseServiceProvider) Register(app *App) {
	dbDriver := os.Getenv("DB_DRIVER")
	config := getDefaultConfig().DbConfig
	dbSession, err := ConnectDB(dbDriver, config)
	if err != nil {
		panic(err)
	}
	app.db = dbSession
	app.dbFunc = func(config *DBConfig) (DBSession, error) {
		if config == nil {
			config = getDefaultConfig().DbConfig
		}
		return ConnectDB(dbDriver, config)
	}
}

func (provider *DatabaseServiceProvider) Boot() {
	//
}
