package routes

import (
	"github.com/lemmego/api/app"
)

func apiRoutes(r app.Router) {
	apiGroup := r.Group("/api")
	{
		apiGroup.Get("/ping", func(c *app.Context) error {
			return app.M{"message": "pong"}
		})
	}
}
