package config

import (
	"github.com/lemmego/lemmego/api/config"
	"github.com/lemmego/lemmego/api/utils"
)

func init() {
	config.SetConfig("app.name", utils.MustEnv("APP_NAME", "Pressebo"))
	config.SetConfig("app.port", utils.MustEnv("APP_PORT", 8080))
	config.SetConfig("app.env", utils.MustEnv("APP_ENV", "development"))
	config.SetConfig("app.debug", utils.MustEnv("APP_DEBUG", "true"))
}
