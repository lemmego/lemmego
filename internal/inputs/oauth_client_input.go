package inputs

import (
	"github.com/lemmego/api/app"
)

type OauthClientInput struct {
	*app.BaseInput
	Name        string `json:"name" in:"form=name"`
	RedirectUri string `json:"redirect_uri" in:"form=redirect_uri"`
}

func (i *OauthClientInput) Validate() error {
	i.Validator.Field("name", i.Name).Required()
	i.Validator.Field("redirect_uri", i.RedirectUri).Required().URL()
	return i.Validator.Validate()
}
