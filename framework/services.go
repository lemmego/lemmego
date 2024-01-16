package framework

type ServiceProvider interface {
	Register(app *App)
	Boot()
}

type BaseServiceProvider struct {
	App *App
}