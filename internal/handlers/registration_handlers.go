package handlers

import (
	"github.com/lemmego/lemmego/api"
	"github.com/lemmego/lemmego/internal/inputs"
	"github.com/lemmego/lemmego/internal/models"
)

func RegistrationStoreHandler(ctx *api.Context) error {
	body := &inputs.RegistrationInput{}
	if err := ctx.Validate(body); err != nil {
		return err
	}

	_, err := ctx.Upload("logo", "images/orgs")

	if err != nil {
		return err
	}

	org := &models.Org{
		OrgUsername: body.OrgUsername,
		OrgName:     body.OrgName,
		OrgEmail:    body.OrgEmail,
	}

	if err := ctx.DB().Create(org).Error; err != nil {
		return err
	}

	user := &models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
		OrgId:     org.ID,
	}

	if err := ctx.DB().Create(user).Error; err != nil {
		return err
	}

	// return ctx.Inertia(200, "Forms/Login", map[string]any{"success": "Registration Successful"})

	return ctx.With("message", "Registration Successful. Please Log In.").Redirect(302, "/login")
}
