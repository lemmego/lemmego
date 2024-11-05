package auth

import "github.com/lemmego/api/app"

func loadAuthRoutes(r app.Router, auth *Auth) {
	r.Get("/login", auth.Guest, func(c *app.Context) error {
		props := map[string]any{}
		message := c.PopSessionString("message")
		if message != "" {
			props["message"] = message
		}
		return c.Inertia(200, "Forms/Login", props)
	})

	r.Get("/register", auth.Guest, func(c *app.Context) error {
		return c.Inertia(200, "Forms/Register", nil)
	})
}
