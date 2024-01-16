package handlers

import (
	"pressebo/framework"
	"pressebo/templates"
	// "pressebo/plugins/auth"
)

func IndexHomeHandler(ctx *framework.Context) error {
	// authUser := ctx.Get("user").(*auth.AuthUser)
	return ctx.Templ(templates.BaseLayout(templates.Hello("John Doe")))
	// return ctx.Render(200, "home.page.tmpl", &fluent.TemplateData{
	// 	StringMap: map[string]string{
	// 		"user": authUser.Username,
	// 	},
	// })
}
