package bootstrap

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/middleware"
)

func LoadHTTPMiddlewares() []app.HTTPMiddleware {
	return []app.HTTPMiddleware{
		middleware.Recoverer(),
		middleware.RequestLogger(),
		middleware.MethodOverride,
	}
}
