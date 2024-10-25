package routes

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/lemmego/internal/handlers"
	"github.com/lemmego/lemmego/internal/plugins"
	"github.com/lemmego/lemmego/internal/plugins/auth"
)

func LoadWebRoutes(r app.Router) {
	//r.Get("/", plugins.Get((*auth.Auth)(nil)).Guard, func(c *app.Context) error {
	//	if user, ok := c.GetSession("user").(*auth.AuthUser); ok {
	//		return c.Inertia(200, "Home/Index", app.M{"user": user})
	//	}
	//	return c.Inertia(200, "Home/Index", nil)
	//})

	r.Get("/oauth/clients/create", handlers.OauthClientCreateHandler)
	r.Post("/oauth/clients", handlers.OauthClientStoreHandler)
	r.Get("/oauth/authorize", handlers.AuthorizeIndexHandler)
	r.Post("/register", plugins.Get(&auth.Auth{}).Guest, handlers.RegistrationStoreHandler)
	r.Post("/login", plugins.Get(&auth.Auth{}).Guest, plugins.Get(&auth.Auth{}).Tenant, handlers.LoginStoreHandler)
}
