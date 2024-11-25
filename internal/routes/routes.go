package routes

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/middleware"
)

func Load() app.RouteCallback {
	// Define your routes here
	return func(r app.Router) {
		r.Use(middleware.Recoverer(), middleware.RequestLogger(), middleware.MethodOverride)
		r.UseBefore(middleware.VerifyCSRF)

		webRoutes(r)
		apiRoutes(r)
		//authRoutes(r)
	}
}
