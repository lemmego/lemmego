package routes

import (
	"github.com/lemmego/api/app"
)

func LoadApiRoutes(r app.Router) {
	apiGroup := r.Group("/api")
	{
		apiGroup.Get("/foo", func(c *app.Context) error {
			return c.Respond(&app.R{
				Payload: app.M{"message": "Hello"},
			})
		})
	}
}
