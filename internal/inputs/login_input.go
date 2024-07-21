package inputs

import (
    "lemmego/api/vee"
)

type LoginInput struct {
    Email string `json:"email" in:"form=email"`
    Password string `json:"password" in:"form=password"`
    OrgUsername string `json:"org_username" in:"form=org_username"`
}

func (i *LoginInput) Validate() error {
	v := vee.New()
    v.Required("email", i.Email)
    v.Required("password", i.Password)
    v.Required("org_username", i.OrgUsername)
	return v.Errors
}
