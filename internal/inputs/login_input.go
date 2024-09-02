package inputs

import (
	"github.com/lemmego/api/app"
)

type LoginInput struct {
	*app.BaseInput
	Email       string `json:"email" in:"form=email"`
	Password    string `json:"password" in:"form=password"`
	OrgUsername string `json:"org_username" in:"form=org_username"`
}

func (i *LoginInput) Validate() error {
	i.Validator.Field("email", i.Email).Required()
	i.Validator.Field("password", i.Password).Required()
	i.Validator.Field("org_username", i.OrgUsername).Required()
	return i.Validator.Validate()
}
