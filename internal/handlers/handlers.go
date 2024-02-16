package handlers

import (
	"pressebo/api"
	"pressebo/templates"
)

func RegisterRoutes(app *api.App) {
	app.Get("/", IndexHomeHandler)
	app.Post("/test", StoreTestHandler)
	app.Get("/foo", func(c *api.Context) error {
		return c.Templ(
			templates.BaseLayout(templates.Hello("Tanmay")),
		)
	})
}
