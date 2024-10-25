package routes

import (
	"github.com/lemmego/api/app"
	baseMw "github.com/lemmego/api/middleware"
	"github.com/lemmego/lemmego/internal/middleware"
)

func Load() func(r app.Router) {
	// Define your routes here
	return func(r app.Router) {
		r.Use(middleware.Logger(), middleware.Recoverer())
		r.UseBefore(baseMw.VerifyCSRF)

		LoadWebRoutes(r)
		LoadApiRoutes(r)
	}
}
