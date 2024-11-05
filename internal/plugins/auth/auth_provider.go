package auth

import (
	_ "embed"
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/db"
	"github.com/lemmego/api/session"
	"github.com/spf13/cobra"
)

//go:embed config.go
var authConfig []byte

type AuthProvider struct {
	*app.ServiceProvider
}

func (p *AuthProvider) Register(a app.AppManager) {
	var dm *db.DatabaseManager
	var sess *session.Session

	if err := a.Service(&dm); err != nil {
		panic(err)
	}

	if err := a.Service(&sess); err != nil {
		panic(err)
	}

	dbc, err := dm.Get()
	if err != nil {
		panic(err)
	}

	a.AddService(New(func(opts *Options) {
		opts.Router = a.Router()
		opts.DB = dbc
		opts.Session = sess
	}))
}

func (p *AuthProvider) Boot(a app.AppManager) {
	var auth *Auth
	if err := a.Service(&auth); err != nil {
		panic(err)
	}

	p.AddRoutes(func(r app.Router) {
		loadAuthRoutes(r, auth)
	})

	p.AddCommands([]app.Command{
		func(a app.AppManager) *cobra.Command {
			return auth.InstallCommand()
		},
	})

	p.Publishes(map[string][]byte{
		"./internal/configs/auth.go": authConfig,
	}, "")
}
