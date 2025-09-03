package bootstrap

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/lemmego/internal/routes"
)

func LoadRoutes() []app.RouteCallback {
	return []app.RouteCallback{
		routes.WebRoutes,
		routes.ApiRoutes,
	}
}
