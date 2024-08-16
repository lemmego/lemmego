package inputs

import (
	"github.com/lemmego/lemmego/api/app"
)

type RegistrationInput struct {
	*app.BaseInput
	Email                string         `json:"email" in:"form=email"`
	Password             string         `json:"password" in:"form=password"`
	PasswordConfirmation string         `json:"password_confirmation" in:"form=password_confirmation"`
	FirstName            string         `json:"first_name" in:"form=first_name"`
	LastName             string         `json:"last_name" in:"form=last_name"`
	Username             string         `json:"username" in:"form=username"`
	Bio                  string         `json:"bio" in:"form=bio"`
	Phone                string         `json:"phone" in:"form=phone"`
	Avatar               *app.FileInput `json:"avatar" in:"form=avatar"`
	OrgName              string         `json:"org_name" in:"form=org_name"`
	OrgEmail             string         `json:"org_email" in:"form=org_email"`
	OrgLogo              *app.FileInput `json:"org_logo" in:"form=org_logo"`
	OrgUsername          string         `json:"org_username" in:"form=org_username"`
}

func (i *RegistrationInput) Validate() error {
	i.Validator.Field("email", i.Email).Required().Unique("users", "email")
	i.Validator.Field("password", i.Password).Required()
	i.Validator.Field("password_confirmation", i.PasswordConfirmation).Required()
	i.Validator.Field("first_name", i.FirstName).Required()
	i.Validator.Field("last_name", i.LastName).Required()
	i.Validator.Field("username", i.Username).Required()
	i.Validator.Field("bio", i.Bio).Required()
	i.Validator.Field("phone", i.Phone).Required()
	i.Validator.Field("avatar", i.Avatar).Required()
	i.Validator.Field("org_name", i.OrgName).Required()
	i.Validator.Field("org_email", i.OrgEmail).Required()
	i.Validator.Field("org_logo", i.OrgLogo).Required()
	i.Validator.Field("org_username", i.OrgUsername).Required()
	return i.Validator.Validate()
}
