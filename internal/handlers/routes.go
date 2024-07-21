package handlers

import (
	"lemmego/api"
)

func Register(r *api.Router) {
	r.Post("/register", RegistrationStoreHandler)
	r.Group("/api", func(r *api.Router) {
		//
	})
	r.Get("/error", func(c *api.Context) error {
		err := c.Pop("error").(string)
		return c.HTML(500, "<html> <body> <pre>"+err+"</body> </html>")
	})
}
