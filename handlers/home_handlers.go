package handlers

import (
	"log"
	"pressebo/framework"
	"pressebo/plugins/auth"
	"pressebo/templates"
	// "pressebo/plugins/auth"
)

func IndexHomeHandler(ctx *framework.Context) error {
	return ctx.HTML(200, `
		<h1>Test Form:</h1>
		<form enctype="multipart/form-data" action="/test" method="POST">
			<input type="text" name="username" placeholder="Username" />
			<input type="password" name="password" placeholder="Password" />
			<input type="file" name="logos[]" placeholder="Logo 1" />
		<input type="file" name="logos[]" placeholder="Logo 2" />
			<input type="submit" value="Submit" />
		</form>
	`)
	// authUser := ctx.Get("user").(*auth.AuthUser)
	return ctx.Templ(templates.BaseLayout(templates.Hello("John Doe")))
	// return ctx.Render(200, "home.page.tmpl", &fluent.TemplateData{
	// 	StringMap: map[string]string{
	// 		"user": authUser.Username,
	// 	},
	// })
}

func StoreTestHandler(ctx *framework.Context) error {
	var req auth.LoginStoreRequest
	ctx.Decode(&req)
	log.Printf("%+v", req)
	return ctx.JSON(200, framework.M{"req": req})
}
