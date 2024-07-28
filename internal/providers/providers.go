package providers

import "lemmego/api"

func Load() []api.ServiceProvider {
	return []api.ServiceProvider{
		&RouteServiceProvider{},
	}
}
