package inputs

import (
	"lemmego/api/vee"
)

type RegistrationInput struct {
	FirstName            string `json:"first_name" in:"form=first_name"`
	LastName             string `json:"last_name" in:"form=last_name"`
	Logo                 string `json:"logo" in:"form=logo"`
	OrgName              string `json:"org_name" in:"form=org_name"`
	OrgEmail             string `json:"org_email" in:"form=org_email"`
	OrgUsername          string `json:"org_username" in:"form=org_username"`
	Email                string `json:"email" in:"form=email"`
	Password             string `json:"password" in:"form=password"`
	PasswordConfirmation string `json:"password_confirmation" in:"form=password_confirmation"`
}

func (i *RegistrationInput) Validate() error {
	v := vee.New()
	v.Required("first_name", i.FirstName)
	v.Required("last_name", i.LastName)
	//v.Required("logo", i.Logo)
	v.Required("org_name", i.OrgName)
	v.Required("org_email", i.OrgEmail)
	v.Required("org_username", i.OrgUsername)

	return v.Validate()
}
