package config

import (
	"lemmego/api"
)

func init() {
	api.SetConfig("app.name", api.MustEnv("APP_NAME", "Pressebo"))
	api.SetConfig("app.port", api.MustEnv("APP_PORT", 8080))
	api.SetConfig("app.env", api.MustEnv("APP_ENV", "development"))
	api.SetConfig("app.debug", api.MustEnv("APP_DEBUG", "true"))
}
