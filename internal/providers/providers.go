package providers

import "github.com/lemmego/api/app"

func Load() []app.Provider {
	return []app.Provider{
		&AppProvider{},
	}
}
