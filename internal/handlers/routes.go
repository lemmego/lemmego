package handlers

import (
	"fmt"
	"lemmego/api"
)

func Routes(r *api.Router) {
	//r.Get("/", func(c *api.Context) error {
	//	return c.Send(200, []byte(c.Query("code")))
	//})
	r.Get("/oauth/clients/create", OauthClientCreateHandler)
	r.Post("/oauth/clients", OauthClientStoreHandler)
	r.Get("/oauth/authorize", AuthorizeIndexHandler)
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
}
