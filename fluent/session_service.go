package fluent

type SessionServiceProvider struct {
	BaseServiceProvider
}

func (provider *SessionServiceProvider) Register(app *App) {
	app.session = NewSessionManager()
}

func (provider *SessionServiceProvider) Boot() {
	//
}
