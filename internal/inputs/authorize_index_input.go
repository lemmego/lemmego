package inputs

import (
	"github.com/lemmego/lemmego/api/app"
	"github.com/lemmego/lemmego/api/vee"
)

type AuthorizeIndexInput struct {
	app.AppManager
	ClientId    string `json:"client_id" in:"query=client_id"`
	State       string `json:"state" in:"query=state"`
	RedirectUri string `json:"redirect_uri" in:"query=redirect_uri"`
	Scope       string `json:"scope" in:"query=scope"`
}

func (i *AuthorizeIndexInput) Validate() error {
	v := vee.New(i.AppManager)
	v.Field("client_id", i.ClientId).Required()
	v.Field("state", i.State).Required()
	v.Field("redirect_uri", i.RedirectUri).Required()
	v.Field("scope", i.Scope).Required()
	return v.Validate()
}
