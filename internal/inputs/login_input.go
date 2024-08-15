package inputs

import (
	"github.com/lemmego/lemmego/api/app"
	"github.com/lemmego/lemmego/api/vee"
)

type LoginInput struct {
	app.AppManager
	Email       string `json:"email" in:"form=email"`
	Password    string `json:"password" in:"form=password"`
	OrgUsername string `json:"org_username" in:"form=org_username"`
}

func (i *LoginInput) Validate() error {
	v := vee.New(i.AppManager)
	v.Field("email", i.Email).Required()
	v.Field("password", i.Password).Required()
	v.Field("org_username", i.OrgUsername).Required()
	return v.Validate()
}
