package handlers

import (
	"fmt"

	"github.com/lemmego/api/app"
	"github.com/lemmego/lemmego/internal/plugins"
	"github.com/lemmego/lemmego/internal/plugins/auth"
)

func Routes(r *app.Router) {

	r.Get("/", plugins.Get(&auth.Auth{}).Guard, func(c *app.Context) error {
		if user, ok := c.GetSession("user").(*auth.AuthUser); ok {
			return c.Inertia(200, "Home/Index", app.M{"user": user})
		}
		return c.Inertia(200, "Home/Index", nil)
	})

	r.Get("/oauth/clients/create", OauthClientCreateHandler)
	r.Post("/oauth/clients", OauthClientStoreHandler)
	r.Get("/oauth/authorize", AuthorizeIndexHandler)
	r.Post("/register", plugins.Get(&auth.Auth{}).Guest, RegistrationStoreHandler)
	r.Post("/login", plugins.Get(&auth.Auth{}).Guest, plugins.Get(&auth.Auth{}).Tenant, LoginStoreHandler)

	apiGroup := r.Group("/api")
	{
		apiGroup.UseBefore(func(c *app.Context) error {
			fmt.Println("before api")
			return c.Next()
		})

		apiGroup.UseAfter(func(c *app.Context) error {
			fmt.Println("after api")
			return c.Next()
		})

		apiGroup.Get("/test3", func(c *app.Context) error {
			c.Text(200, []byte("test3"))
			return c.Next()
		})

		v1Group := apiGroup.Group("/v1")
		{
			v1Group.UseBefore(func(c *app.Context) error {
				fmt.Println("before v1")
				return c.Next()
			})
			v1Group.UseAfter(func(c *app.Context) error {
				fmt.Println("after v1")
				return c.Next()
			})

			v1Group.Get("/test1", func(c *app.Context) error {
				fmt.Println("inside test1")
				c.Text(200, []byte("test1"))
				return c.Next()
			})

			v1Group.Get("/test2", func(c *app.Context) error {
				c.Text(200, []byte("test2"))
				return c.Next()
			}).UseBefore(func(c *app.Context) error {
				println("before test2")
				return c.Next()
			}).UseAfter(func(c *app.Context) error {
				fmt.Println("after test2")
				return c.Next()
			})
		}
	}

}
