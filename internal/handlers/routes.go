package handlers

import (
	"fmt"
	"lemmego/api"
)

func Register(r *api.Router) {
	r.Post("/register", RegistrationStoreHandler)
	r.Group("/api", func(r *api.Router) {
		r.Group("/v1", func(r *api.Router) {
			r.Get("/test1", func(c *api.Context) error {
				return c.Send(200, []byte("Hello from test1"))
			}).UseBefore(func(next api.Handler) api.Handler {
				return func(c *api.Context) error {
					fmt.Println("Executing test1 middleware")
					return next(c)
				}
			})
			r.Get("/test2", func(c *api.Context) error {
				return c.Send(200, []byte("Hello from test2"))
			})
		}).UseBefore(func(next api.Handler) api.Handler {
			return func(c *api.Context) error {
				fmt.Println("Executing v1 middleware")
				return next(c)
			}
		})
	}).UseBefore(func(next api.Handler) api.Handler {
		return func(c *api.Context) error {
			fmt.Println("Executing api middleware")
			return next(c)
		}
	})

	r.Get("/error", func(c *api.Context) error {
		err := c.Pop("error").(string)
		return c.HTML(500, "<html> <body> <pre>"+err+"</body> </html>")
	})
}
