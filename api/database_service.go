package api

type DatabaseServiceProvider struct {
	BaseServiceProvider
}

func (provider *DatabaseServiceProvider) Register(app *App) {
	dbConfig := &DBConfig{
		Driver:   Config("db.driver").(string),
		Host:     Config("db.host").(string),
		Port:     Config("db.port").(int),
		Database: Config("db.database").(string),
		User:     Config("db.username").(string),
		Password: Config("db.password").(string),
	}
	dbSession, err := ConnectDB(dbConfig.Driver, dbConfig)
	if err != nil {
		panic(err)
	}
	app.db = dbSession
	app.dbFunc = func(config *DBConfig) (DBSession, error) {
		if config == nil {
			config = dbConfig
		}
		return ConnectDB(dbConfig.Driver, config)
	}
}

func (provider *DatabaseServiceProvider) Boot() {
	//
}
