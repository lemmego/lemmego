package bootstrap

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/providers/fs"
	"github.com/lemmego/api/providers/session"
	"github.com/lemmego/auth"
	"github.com/lemmego/gormconnector"
	"github.com/lemmego/inertia"
)

func LoadProviders() []app.Provider {
	return []app.Provider{
		&fs.Provider{},
		&session.Provider{},
		&inertia.Provider{},
		&gormconnector.Provider{
			UseGPA: true,
		},
		&auth.Provider{
			Opts: &auth.Opts{
				DisableSession: true,
				JwtSecret:      "a-long-long-secret",
			},
		},
	}
}
