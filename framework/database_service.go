package framework

type DatabaseServiceProvider struct {
	BaseServiceProvider
}

func (provider *DatabaseServiceProvider) Register(app *App) {
	dbSession, err := ConnectDB("postgres", getDefaultConfig().DbConfig)
	if err != nil {
		panic(err)
	}
	app.db = dbSession
	app.dbFunc = func(config *DBConfig) (DBSession, error) {
		if config == nil {
			config = getDefaultConfig().DbConfig
		}
		return ConnectDB("postgres", config)
	}
}

func (provider *DatabaseServiceProvider) Boot() {
	//
}
