package inputs

import (
	"github.com/lemmego/lemmego/api/app"
	"github.com/lemmego/lemmego/api/vee"
)

type OauthClientInput struct {
	app.AppManager
	Name        string `json:"name" in:"form=name"`
	RedirectUri string `json:"redirect_uri" in:"form=redirect_uri"`
}

func (i *OauthClientInput) Validate() error {
	v := vee.New(i.AppManager)
	v.Field("name", i.Name).Required()
	v.Field("redirect_uri", i.RedirectUri).Required().URL()
	return v.Validate()
}
