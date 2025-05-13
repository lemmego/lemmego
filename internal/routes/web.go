package routes

import (
	"github.com/lemmego/api/app"
)

func webRoutes(r app.Router) {
	r.Get("/{$}", func(c *app.Context) error {
		// return c.Inertia("IndexVue", nil)
		// return c.Inertia("IndexReact", nil)
		return c.Render("index.page.gohtml", nil)
	})
}

func init() {
	app.RegisterRoutes(webRoutes)
}
