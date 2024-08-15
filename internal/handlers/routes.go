package handlers

import (
	"fmt"

	"github.com/lemmego/lemmego/api/app"
)

func Routes(r *app.Router) {
	//r.Get("/", func(c *api.Context) error {
	//	return c.Send(200, []byte(c.Query("code")))
	//})

	r.Get("/oauth/clients/create", OauthClientCreateHandler)
	r.Post("/oauth/clients", OauthClientStoreHandler)
	r.Get("/oauth/authorize", AuthorizeIndexHandler)
	r.Post("/register", RegistrationStoreHandler)
	//r.Get("/api/v1Group/test1", func(c *api.Context) error {
	//	fmt.Println("inside test1")
	//	return c.Send(200, []byte("test1"))
	//}).UseBefore(func(next api.Handler) api.Handler {
	//	return func(c *api.Context) error {
	//		fmt.Println("before test1")
	//		return next(c)
	//	}
	//}).UseAfter(func(next api.Handler) api.Handler {
	//	return func(c *api.Context) error {
	//		fmt.Println("after test1")
	//		return next(c)
	//	}
	//})

	apiGroup := r.Group("/api")

	apiGroup.UseBefore(func(c *app.Context) error {
		fmt.Println("before api")
		return c.Next()
	})

	apiGroup.UseAfter(func(c *app.Context) error {
		fmt.Println("after api")
		return c.Next()
	})

	apiGroup.Get("/test3", func(c *app.Context) error {
		c.Send(200, []byte("test3"))
		return c.Next()
	})

	v1Group := apiGroup.Group("/v1")
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
		c.Send(200, []byte("test1"))
		return c.Next()
	})

	v1Group.Get("/test2", func(c *app.Context) error {
		c.Send(200, []byte("test2"))
		return c.Next()
	}).UseBefore(func(c *app.Context) error {
		println("before test2")
		return c.Next()
	}).UseAfter(func(c *app.Context) error {
		fmt.Println("after test2")
		return c.Next()
	})

}
