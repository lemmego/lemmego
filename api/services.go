package api

type ServiceProvider interface {
	Register(app *App)
	Boot()
}

type BaseServiceProvider struct {
	App *App
}

func (p *BaseServiceProvider) Publishes() {
	//
}
