package providers

import "github.com/lemmego/api/app"

func Load() []app.ServiceProvider {
	return []app.ServiceProvider{
		&RouteServiceProvider{},
	}
}
