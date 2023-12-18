package fluent

type DatabaseServiceProvider struct {
	BaseServiceProvider
}

func (provider *DatabaseServiceProvider) Register(app *App) {
	app.Container.Singleton(func() DBSession {
		sess, err := ConnectDB("postgres", getDefaultConfig().DbConfig)
		if err != nil {
			panic(err)
		}
		return sess
	})
}

func (provider *DatabaseServiceProvider) Boot() {
	//
}
