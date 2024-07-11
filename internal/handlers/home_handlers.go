package handlers

import (
	"errors"
	"lemmego/api"
	"lemmego/templates"
	"net/http"
)

func IndexHomeHandler(ctx *api.Context) error {
	ctx.Put("foo", "Bar")
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"foo": ctx.Pop("foo"),
	})
	err := ctx.Inertia("Home/Welcome", map[string]any{
		"name": "John Doe",
	})
	if err != nil {
		return ctx.Unauthorized(err)
	}
	return ctx.Unauthorized(errors.New("Unauthorized"))
	return ctx.HTML(200, `
		<h1>Test Form:</h1>
		<form enctype="multipart/form-data" action="/test" method="POST">
			<input type="text" name="username" placeholder="Username" />
			<input type="password" name="password" placeholder="Password" />
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

func StoreTestHandler(ctx *api.Context) error {
	input := &TestInput{}
	if validated, err := ctx.ParseAndValidate(input); err != nil {
		return ctx.JSON(400, api.M{"errors": err})
	} else {
		input = validated.(*TestInput)
	}

	return ctx.JSON(200, api.M{"input": input})
}
