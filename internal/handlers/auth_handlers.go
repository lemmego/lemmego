package handlers

import (
	"fmt"
	"lemmego/api"
	"lemmego/internal/inputs"
)

func AuthIndexHandler(ctx *api.Context) error {
	return nil
}

func AuthCreateHandler(ctx *api.Context) error {
	return nil
}

func AuthShowHandler(ctx *api.Context) error {
	return nil
}

func AuthStoreHandler(ctx *api.Context) error {
	if body, err := ctx.ParseAndValidate(&inputs.Registration{}); err != nil {
		fmt.Println("errors in handler", err.Error())
		return err
	} else {
		input := body.(*inputs.Registration)

		ctx.App().Inertia().ShareProp("input", input)
		ctx.App().Inertia().Back(ctx.ResponseWriter(), ctx.Request())
	}

	return nil
}

func AuthEditHandler(ctx *api.Context) error {
	return nil
}

func AuthUpdateHandler(ctx *api.Context) error {
	return nil
}

func AuthDeleteHandler(ctx *api.Context) error {
	return nil
}
