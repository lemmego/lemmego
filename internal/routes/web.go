package routes

import (
	"github.com/lemmego/api/app"
)

func webRoutes(r app.Router) {
	r.Get("/{$}", func(c *app.Context) error {
		return c.Inertia("Index", nil)
		return c.Render("index.page.gohtml", nil)
	})
}
