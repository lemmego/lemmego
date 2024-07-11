package handlers

import (
	"lemmego/api"
	"lemmego/templates"
	"log/slog"
)

func LoginIndexHandler(c *api.Context) error {
	return c.Templ(templates.BaseLayout(templates.Login()))
}

func RegistrationIndexHandler(c *api.Context) error {
	return c.Inertia("Forms/Register", map[string]any{"foo": "bar"})
	// return c.App().i.Render(c.ResponseWriter(), c.Request(), "Forms/Register", nil)
	// return c.Templ(templates.BaseLayout(templates.Register()))
}

func LoginStoreHandler(c *api.Context) error {
	return nil
}

func RegistrationStoreHandler(c *api.Context) error {
	body := c.GetBody()
	slog.Info("parsed", "body", body)
	return c.Back()
}
