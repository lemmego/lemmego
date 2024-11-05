package handlers

import (
	"errors"
	"log/slog"

	"github.com/lemmego/api/app"
	"github.com/lemmego/api/db"
	"github.com/lemmego/api/shared"
	"github.com/lemmego/lemmego/internal/inputs"
	"github.com/lemmego/lemmego/internal/models"
	"github.com/lemmego/lemmego/internal/plugins"
	"github.com/lemmego/lemmego/internal/plugins/auth"
)

func LoginStoreHandler(c *app.Context) error {
	var dm *db.DatabaseManager
	if err := c.App().Service(&dm); err != nil {
		return err
	}

	credErrors := shared.ValidationErrors{
		"password": []string{"Invalid credentials"},
		"email":    []string{"Invalid credentials"},
	}
	body := &inputs.LoginInput{}
	if err := c.Validate(body); err != nil {
		return err
	}
	user := &models.User{Email: body.Email}

	user.OrgId = c.Get("org_id").(uint)

	conn, err := dm.Get()
	if err != nil {
		return err
	}

	if err := conn.DB().Debug().Where(user).First(user).Error; err != nil {
		return err
	}

	if user.ID == 0 {
		return credErrors
	}

	slog.Info("User logged in", "user", user)

	authPlugin := plugins.Get(&auth.Auth{})

	if _, err := authPlugin.Login(
		c.Request().Context(),
		&auth.CredUser{ID: user.ID, Username: user.Email, Password: user.Password},
		body.Email,
		body.Password,
	); err != nil {
		println("error logging in", err)
		if errors.Is(err, auth.ErrInvalidCreds) {
			return credErrors
		}
		return err
	}
	return c.Redirect(302, "/")
}

func LoginDeleteHandler(ctx *app.Context) error {
	return nil
}
