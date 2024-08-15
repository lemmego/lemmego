package handlers

import (
	"fmt"

	"github.com/lemmego/lemmego/api/app"
	"github.com/lemmego/lemmego/internal/inputs"
	"github.com/lemmego/lemmego/internal/models"
)

func RegistrationStoreHandler(ctx *app.Context) error {
	body := &inputs.RegistrationInput{}
	if err := ctx.Validate(body); err != nil {
		return err
	}

	_, err := ctx.Upload("org_logo", "images/orgs")

	if err != nil {
		return fmt.Errorf("could not upload org_logo: %w", err)
	}

	_, err = ctx.Upload("avatar", "images/orgs")

	if err != nil {
		return fmt.Errorf("could not upload avatar: %w", err)
	}

	org := &models.Org{
		OrgUsername: body.OrgUsername,
		OrgName:     body.OrgName,
		OrgEmail:    body.OrgEmail,
	}

	if ctx.HasFile("org_logo") {
		org.OrgLogo = "images/orgs/" + body.OrgLogo.Filename()
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
		Bio:       body.Bio,
		Phone:     body.Phone,
		Username:  body.Username,
	}

	if ctx.HasFile("avatar") {
		user.Avatar = "images/orgs/" + body.Avatar.Filename()
	}

	if err := ctx.DB().Create(user).Error; err != nil {
		return err
	}

	// return ctx.Inertia(200, "Forms/Login", map[string]any{"success": "Registration Successful"})

	return ctx.With("message", "Registration Successful. Please Log In.").Redirect(302, "/login")
}
