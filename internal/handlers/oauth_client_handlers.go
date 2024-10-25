package handlers

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/db"
	"github.com/lemmego/lemmego/internal/inputs"
	"github.com/lemmego/lemmego/internal/models"
	"github.com/lucsky/cuid"
)

func OauthClientIndexHandler(ctx *app.Context) error {
	return nil
}

func OauthClientCreateHandler(ctx *app.Context) error {
	return ctx.Inertia(200, "Forms/OauthClient", nil)
}

func OauthClientShowHandler(ctx *app.Context) error {
	return nil
}

func OauthClientStoreHandler(ctx *app.Context) error {
	var dm *db.DatabaseManager
	if err := ctx.App().Service(&dm); err != nil {
		return err
	}

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

	conn, err := dm.Get()
	if err != nil {
		return err
	}

	if err := conn.DB().Create(&client).Error; err != nil {
		return err
	}

	return ctx.JSON(201, app.M{"message": "Client Created", "data": client})
}

func OauthClientEditHandler(ctx *app.Context) error {
	return nil
}

func OauthClientUpdateHandler(ctx *app.Context) error {
	return nil
}

func OauthClientDeleteHandler(ctx *app.Context) error {
	return nil
}
