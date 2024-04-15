package handlers

import (
	"pressebo/api"
	"pressebo/templates"
)

func LoginIndexHandler(c *api.Context) error {
	return c.Templ(templates.BaseLayout(templates.Login()))
}

func RegistrationIndexHandler(c *api.Context) error {
	return c.Templ(templates.BaseLayout(templates.Register()))
}

func LoginStoreHandler(c *api.Context) error {
	return nil
}

func RegistrationStoreHandler(c *api.Context) error {
	return nil
}
