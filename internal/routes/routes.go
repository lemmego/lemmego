package routes

import (
	"github.com/lemmego/api/app"
	baseMw "github.com/lemmego/api/middleware"
	"github.com/lemmego/lemmego/internal/middleware"
)

func Load() func(r app.Router) {
	return func(r app.Router) {
		r.Use(middleware.Logger(), middleware.Recoverer())
		r.UseBefore(baseMw.VerifyCSRF)

		LoadWebRoutes(r)
		LoadApiRoutes(r)

		r.Get("/error", func(c *app.Context) error {
			err := c.PopSession("error").(string)
			return c.HTML(500, []byte("<html><body><code>"+err+"</code></body></html>"))
		})
	}
}
