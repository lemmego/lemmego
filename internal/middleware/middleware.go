package middleware

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/middleware"
)

type MiddlewareLoader struct {
	GlobalHTTP []app.HTTPMiddleware
	Global     []app.Middleware
	Group      map[string]app.Middleware
}

func Load() MiddlewareLoader {
	return MiddlewareLoader{
		GlobalHTTP: []app.HTTPMiddleware{
			middleware.Logger(),
			middleware.Recoverer(),
		},
	}
}
