package fluent

type SessionServiceProvider struct {
	BaseServiceProvider
}

func (provider *SessionServiceProvider) Register(app *App) {
	app.session = NewSessionManager()
	app.Container.NamedSingleton("session", func() *Session {
		return app.session
	})
}

func (provider *SessionServiceProvider) Boot() {
	//
}
