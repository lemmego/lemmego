package handlers

import (
	"pressebo/api"
	// "pressebo/plugins/auth"
)

func RegisterRoutes(app *api.App) {
	// auth := app.Plugin(auth.Namespace).(*auth.AuthPlugin)
	app.Get("/", IndexHomeHandler)
	app.Post("/test", StoreTestHandler)
	app.Get("/foo", func(c *api.Context) error {
		return c.JSON(200, api.M{"foo": "bar"})
	})
}
