package handlers

import (
	"errors"
	"strconv"

	"github.com/lemmego/lemmego/api/app"
	"github.com/lemmego/lemmego/api/logger"
	"github.com/lemmego/lemmego/api/shared"
	"github.com/lemmego/lemmego/internal/inputs"
	"github.com/lemmego/lemmego/internal/models"
	"github.com/lemmego/lemmego/internal/plugins"
	"github.com/lemmego/lemmego/internal/plugins/auth"
)

func LoginStoreHandler(c *app.Context) error {
	credErrors := shared.ValidationErrors{
		"password": []string{"Invalid credentials"},
		"email":    []string{"Invalid credentials"},
	}
	body := &inputs.LoginInput{}
	if err := c.Validate(body); err != nil {
		return err
	}
	user := &models.User{Email: body.Email}

	c.DB().Where(user).First(user)

	if user.ID == 0 {
		return credErrors
	}

	logger.D().Info("User logged in", "user", user)

	authPlugin := plugins.Get(&auth.AuthPlugin{})

	if _, err := authPlugin.Login(
		c.Request().Context(),
		&auth.CredUser{ID: strconv.Itoa(int(user.ID)), Username: user.Email, Password: user.Password},
		body.Email,
		body.Password,
	); err != nil {
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
