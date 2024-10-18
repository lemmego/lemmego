package auth

import "github.com/lemmego/api/app"

type AuthProvider struct {
	*app.ServiceProvider
}

func (p *AuthProvider) Register(a app.AppManager) {
	//TODO implement me
	//auth := New()
	//p.AddRoutes(func(r app.Router) {
	//
	//})
}

func (p *AuthProvider) Boot(a app.AppManager) {
	//TODO implement me
}
