package routes

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/res"
)

func WebRoutes(a app.App) {
	r := a.Router()
	r.Get("/{$}", func(c app.Context) error {
		//return inertia.Respond(c, "IndexReact", nil) // Requires lemmego/inertia plugin
		//return inertia.Respond(c, "IndexVue", nil) // Requires lemmego/inertia plugin
		//return templ.Respond(c, templates.BaseLayout(templates.Index())) // Requires lemmego/templ plugin
		return c.Render(res.NewTemplate(c, "index.page.gohtml"))
	})
}
