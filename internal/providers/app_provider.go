package providers

import (
	"github.com/lemmego/api/app"
)

type AppProvider struct {
	*app.ServiceProvider
}

func (p *AppProvider) Register(a app.AppManager) {
	//TODO implement me
}

func (p *AppProvider) Boot(a app.AppManager) {
	//TODO implement me
}
