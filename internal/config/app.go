package config

import (
	"github.com/lemmego/lemmego/api"
	"github.com/lemmego/lemmego/api/config"
)

func init() {
	config.SetConfig("app.name", api.MustEnv("APP_NAME", "Pressebo"))
	config.SetConfig("app.port", api.MustEnv("APP_PORT", 8080))
	config.SetConfig("app.env", api.MustEnv("APP_ENV", "development"))
	config.SetConfig("app.debug", api.MustEnv("APP_DEBUG", "true"))
}
