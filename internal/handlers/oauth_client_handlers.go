package handlers

import (
	"github.com/lucsky/cuid"
	"lemmego/api"
	"lemmego/internal/inputs"
	"lemmego/internal/models"
)

func OauthClientIndexHandler(ctx *api.Context) error {
	return nil
}

func OauthClientCreateHandler(ctx *api.Context) error {
	return ctx.Inertia(200, "Forms/OauthClient", nil)
}

func OauthClientShowHandler(ctx *api.Context) error {
	return nil
}

func OauthClientStoreHandler(ctx *api.Context) error {
	body := &inputs.OauthClientInput{}
	if err := ctx.Validate(body); err != nil {
		return err
	}

	client := models.OauthClient{
		ID:          cuid.New(),
		Secret:      "abcdefghijkl",
		RedirectUri: body.RedirectUri,
		Name:        body.Name,
	}

	if err := ctx.DB().Create(&client).Error; err != nil {
		return err
	}

	return ctx.JSON(201, api.M{"message": "Client Created", "data": client})
}

func OauthClientEditHandler(ctx *api.Context) error {
	return nil
}

func OauthClientUpdateHandler(ctx *api.Context) error {
	return nil
}

func OauthClientDeleteHandler(ctx *api.Context) error {
	return nil
}
