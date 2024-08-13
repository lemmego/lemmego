package providers

import "github.com/lemmego/lemmego/api"

func Load() []api.ServiceProvider {
	return []api.ServiceProvider{
		&RouteServiceProvider{},
	}
}
