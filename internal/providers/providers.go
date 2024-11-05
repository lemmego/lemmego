package providers

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/lemmego/internal/plugins/auth"
)

func Load() []app.Provider {
	return []app.Provider{
		&AppProvider{},
		&auth.AuthProvider{},
	}
}
