package inputs

import (
	"github.com/ggicci/httpin"
	"github.com/lemmego/lemmego/api"
	"github.com/lemmego/lemmego/api/vee"
)

type RegistrationInput struct {
	api.AppManager
	Email                string       `json:"email" in:"form=email"`
	Password             string       `json:"password" in:"form=password"`
	PasswordConfirmation string       `json:"password_confirmation" in:"form=password_confirmation"`
	FirstName            string       `json:"first_name" in:"form=first_name"`
	LastName             string       `json:"last_name" in:"form=last_name"`
	Username             string       `json:"username" in:"form=username"`
	Bio                  string       `json:"bio" in:"form=bio"`
	Phone                string       `json:"phone" in:"form=phone"`
	Avatar               *httpin.File `json:"avatar" in:"form=avatar"`
	OrgName              string       `json:"org_name" in:"form=org_name"`
	OrgEmail             string       `json:"org_email" in:"form=org_email"`
	OrgLogo              *httpin.File `json:"org_logo" in:"form=org_logo"`
	OrgUsername          string       `json:"org_username" in:"form=org_username"`
}

func (i *RegistrationInput) Validate() error {
	v := vee.New()
	v.Field("email", i.Email).Required().Unique(i.DB(), "users", "email")
	v.Field("password", i.Password).Required()
	v.Field("password_confirmation", i.PasswordConfirmation).Required()
	v.Field("first_name", i.FirstName).Required()
	v.Field("last_name", i.LastName).Required()
	v.Field("username", i.Username).Required()
	v.Field("bio", i.Bio).Required()
	v.Field("phone", i.Phone).Required()
	// v.Field("avatar", i.Avatar).Required()
	v.Field("org_name", i.OrgName).Required()
	v.Field("org_email", i.OrgEmail).Required()
	// v.Field("org_logo", i.OrgLogo).Required()
	v.Field("org_username", i.OrgUsername).Required()
	return v.Validate()
}
