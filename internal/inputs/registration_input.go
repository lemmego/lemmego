package inputs

import (
    "github.com/lemmego/lemmego/api/vee"
)

type RegistrationInput struct {
    FirstName string `json:"first_name" in:"form=first_name"`
    LastName string `json:"last_name" in:"form=last_name"`
    Logo string `json:"logo" in:"form=logo"`
    OrgName string `json:"org_name" in:"form=org_name"`
    OrgEmail string `json:"org_email" in:"form=org_email"`
    OrgUsername string `json:"org_username" in:"form=org_username"`
    Email string `json:"email" in:"form=email"`
    Password string `json:"password" in:"form=password"`
    PasswordConfirmation string `json:"password_confirmation" in:"form=password_confirmation"`
}

func (i *RegistrationInput) Validate() error {
	v := vee.New()
    v.Field("first_name", i.FirstName).Required()
    v.Field("last_name", i.LastName).Required()
    v.Field("org_name", i.OrgName).Required()
    v.Field("org_email", i.OrgEmail).Required()
    v.Field("org_username", i.OrgUsername).Required()
    v.Field("email", i.Email).Required()
    v.Field("password", i.Password).Required()
    v.Field("password_confirmation", i.PasswordConfirmation).Required()
	return v.Validate()
}
