package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lemmego/api/app"
)

func Recoverer() app.HTTPMiddleware {
	return middleware.Recoverer
}
