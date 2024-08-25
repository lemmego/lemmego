package config

import (
	"github.com/lemmego/lemmego/api/config"
)

func init() {
	config.Set("app.name", config.MustEnv("APP_NAME", "Pressebo"))
	config.Set("app.port", config.MustEnv("APP_PORT", 8080))
	config.Set("app.env", config.MustEnv("APP_ENV", "development"))
	config.Set("app.debug", config.MustEnv("APP_DEBUG", "true"))
}
