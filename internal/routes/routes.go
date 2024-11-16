package routes

import (
	"github.com/lemmego/api/app"
	mw "github.com/lemmego/api/middleware"
	"github.com/lemmego/lemmego/internal/middleware"
)

func Load() app.RouteCallback {
	// Define your routes here
	return func(r app.Router) {
		r.Use(mw.Recoverer(), middleware.LogRequest(), mw.MethodOverride)
		r.UseBefore(mw.VerifyCSRF)

		webRoutes(r)
		apiRoutes(r)
		//authRoutes(r)
	}
}
