package providers

import (
	"github.com/lemmego/api/app"
)

type GithubServiceProvider struct {
	app.BaseServiceProvider
}

func (provider *GithubServiceProvider) Register(app *app.App) {
	// TODO: Implement
}

func (provider *GithubServiceProvider) Boot() {
	// TODO: Implement
}
