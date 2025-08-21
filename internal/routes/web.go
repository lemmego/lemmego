package routes

import (
	"fmt"
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/db"
	"github.com/lemmego/lemmego/internal/plugins/inertia"
	"github.com/lemmego/lemmego/internal/plugins/templ"
	"github.com/lemmego/lemmego/templates"
)

func WebRoutes(r app.Router) {
	r.Get("/{$}", func(c app.Context) error {
		// return c.Inertia("IndexVue", nil)
		// return c.Inertia("IndexReact", nil)
		// return c.Render("index.page.gohtml", nil)
		// return c.Render(res.NewTemplate(c, "index.page.gohtml"))
		return inertia.Respond(c, "IndexReact", nil)
		return c.Render(templ.New(c, templates.BaseLayout(templates.Index())))
	})

	r.Get("/foo", func(c app.Context) error {
		fmt.Println(db.SqlProvider().ProviderInfo())
		return app.M{"foo": "baz"}
	})
}
