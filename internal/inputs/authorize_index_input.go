package inputs

import (
	"github.com/lemmego/api/app"
)

type AuthorizeIndexInput struct {
	*app.BaseInput
	ClientId    string `json:"client_id" in:"query=client_id"`
	State       string `json:"state" in:"query=state"`
	RedirectUri string `json:"redirect_uri" in:"query=redirect_uri"`
	Scope       string `json:"scope" in:"query=scope"`
}

func (i *AuthorizeIndexInput) Validate() error {
	i.Validator.Field("client_id", i.ClientId).Required()
	i.Validator.Field("state", i.State).Required()
	i.Validator.Field("redirect_uri", i.RedirectUri).Required()
	i.Validator.Field("scope", i.Scope).Required()
	return i.Validator.Validate()
}
