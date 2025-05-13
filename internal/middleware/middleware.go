package middleware

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/middleware"
)

func init() {
	app.RegisterHTTPMiddleware(middleware.Recoverer(), middleware.RequestLogger(), middleware.MethodOverride)
	app.RegisterMiddleware(middleware.VerifyCSRF)
}
