package bootstrap

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/config"
	"github.com/lemmego/api/providers/fs"
	"github.com/lemmego/api/providers/session"
	"github.com/lemmego/auth"
	"github.com/lemmego/gormconnector"
	"github.com/lemmego/inertia"
	"github.com/lemmego/queue"
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
				JwtSecret:      config.MustEnv("JWT_SECRET", config.MustEnv("APP_KEY", "")),
			},
		},
		&queue.Provider{},
	}
}
