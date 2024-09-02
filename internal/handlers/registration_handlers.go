package handlers

import (
	"fmt"

	"github.com/lemmego/api/app"
	"github.com/lemmego/api/db"
	"github.com/lemmego/api/utils"
	"github.com/lemmego/lemmego/internal/inputs"
	"github.com/lemmego/lemmego/internal/models"
)

func RegistrationStoreHandler(c *app.Context) error {
	body := &inputs.RegistrationInput{}
	if err := c.Validate(body); err != nil {
		return err
	}

	password, err := utils.Bcrypt(body.Password)

	if err != nil {
		return err
	}

	org := &models.Org{
		OrgUsername: body.OrgUsername,
		OrgName:     body.OrgName,
		OrgEmail:    body.OrgEmail,
	}

	user := &models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  password,
		Bio:       body.Bio,
		Phone:     body.Phone,
		Username:  body.Username,
	}

	if c.HasFile("org_logo") {
		_, err := c.Upload("org_logo", "images/orgs")

		if err != nil {
			return fmt.Errorf("could not upload org_logo: %w", err)
		}
		org.OrgLogo = "images/orgs/" + body.OrgLogo.Filename()
	}

	if c.HasFile("avatar") {
		_, err := c.Upload("avatar", "images/avatars")

		if err != nil {
			return fmt.Errorf("could not upload avatar: %w", err)
		}
		user.Avatar = "images/avatars/" + body.Avatar.Filename()
	}

	if err := db.Get().Create(org).Error; err != nil {
		return err
	} else {
		user.OrgId = org.ID
	}

	if err := db.Get().Create(user).Error; err != nil {
		return err
	}

	return c.With("message", "Registration Successful. Please Log In.").Redirect(302, "/login")
}
