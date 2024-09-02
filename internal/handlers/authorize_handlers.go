package handlers

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/lemmego/internal/inputs"
)

func AuthorizeIndexHandler(ctx *app.Context) error {
	input := &inputs.AuthorizeIndexInput{}
	if err := ctx.Validate(input); err != nil {
		return err
	}
	return ctx.Inertia(200, "Forms/OauthAuthorize", nil)
}

func AuthorizeCreateHandler(ctx *app.Context) error {
	return nil
}

func AuthorizeShowHandler(ctx *app.Context) error {
	return nil
}

func AuthorizeStoreHandler(ctx *app.Context) error {
	return nil
}

func AuthorizeEditHandler(ctx *app.Context) error {
	return nil
}

func AuthorizeUpdateHandler(ctx *app.Context) error {
	return nil
}

func AuthorizeDeleteHandler(ctx *app.Context) error {
	return nil
}
