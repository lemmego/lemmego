package api

import "pressebo/api/db"

type DatabaseServiceProvider struct {
	BaseServiceProvider
}

func (provider *DatabaseServiceProvider) Register(app *App) {
	dbConfig := &db.DBConfig{
		Driver:   Config("db.driver").(string),
		Host:     Config("db.host").(string),
		Port:     Config("db.port").(int),
		Database: Config("db.database").(string),
		User:     Config("db.username").(string),
		Password: Config("db.password").(string),
	}

	dbSession, err := db.NewConnection(dbConfig).WithForceCreateDb().Open()
	if err != nil {
		panic(err)
	}

	app.db = dbSession
	app.dbFunc = func(config *db.DBConfig) (db.DBSession, error) {
		if config == nil {
			config = dbConfig
		}
		return db.NewConnection(dbConfig).WithForceCreateDb().Open()
	}

}

func (provider *DatabaseServiceProvider) Boot() {
	//
}
