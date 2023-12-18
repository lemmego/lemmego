package fluent

type ServiceProvider interface {
	Register(app *App)
	Boot()
}

type BaseServiceProvider struct {
	App *App
}
