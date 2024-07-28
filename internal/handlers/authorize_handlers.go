package handlers

import (
	"lemmego/api"
	"lemmego/internal/inputs"
)

func AuthorizeIndexHandler(ctx *api.Context) error {
	input := &inputs.AuthorizeIndexInput{}
	if err := ctx.Validate(input); err != nil {
		return err
	}
	return ctx.Inertia(200, "Forms/OauthAuthorize", nil)
}

func AuthorizeCreateHandler(ctx *api.Context) error {
	return nil
}

func AuthorizeShowHandler(ctx *api.Context) error {
	return nil
}

func AuthorizeStoreHandler(ctx *api.Context) error {
	return nil
}

func AuthorizeEditHandler(ctx *api.Context) error {
	return nil
}

func AuthorizeUpdateHandler(ctx *api.Context) error {
	return nil
}

func AuthorizeDeleteHandler(ctx *api.Context) error {
	return nil
}
