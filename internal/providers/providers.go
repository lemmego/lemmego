package providers

import "github.com/lemmego/lemmego/api/app"

func Load() []app.ServiceProvider {
	return []app.ServiceProvider{
		&RouteServiceProvider{},
	}
}
