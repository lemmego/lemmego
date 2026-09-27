package bootstrap

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/config"
	"github.com/lemmego/api/providers/fs"
	"github.com/lemmego/api/providers/session"
	"github.com/lemmego/auth"
	"github.com/lemmego/inertia"
	"github.com/lemmego/lemmego/internal/repos"
	"github.com/lemmego/ormconnector"
	"github.com/lemmego/queue"
)

func LoadProviders() []app.Provider {
	return []app.Provider{
		&fs.Provider{},
		&session.Provider{},
		&inertia.Provider{},
		&ormconnector.Provider{
			UseGPA: true,
		},
		&auth.Provider{
			Opts: &auth.Opts{
				// A credential carries an id, never a copy of the user. This
				// is the one place that turns that id back into a row, so
				// every request resolves the same concrete type no matter
				// which credential arrived — a session cookie, a bearer
				// token, or an OAuth2 access token. Handlers read it with
				// auth.UserAs[*models.User](c).
				//
				// It costs one indexed lookup per authenticated request, and
				// buys the thing a token cannot give you: the row as it is
				// now, not as it was when the token was minted.
				//
				// The loader takes app.Context because LoadProviders runs
				// before the container exists — there is no app.App to close
				// over here. c.App() resolves it per request.
				UserLoader: func(c app.Context, id string) (any, error) {
					return repos.User().FindByID(c.RequestContext(), id)
				},
				// This template serves pages, so it keeps a session: logging
				// out actually revokes, because the server holds the record.
				// JwtSecret is still set, so /api/* takes bearer tokens too.
				DisableSession: false,
				JwtSecret:      config.MustEnv("JWT_SECRET", config.MustEnv("APP_KEY", "")),
			},
		},
		&queue.Provider{},
	}
}
