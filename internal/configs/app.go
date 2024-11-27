package configs

import "github.com/lemmego/api/config"

var app = config.M{
	"name":  config.MustEnv("APP_NAME", "Lemmego"),
	"port":  config.MustEnv("APP_PORT", 8080),
	"env":   config.MustEnv("APP_ENV", "development"),
	"debug": config.MustEnv("APP_DEBUG", false),
}
