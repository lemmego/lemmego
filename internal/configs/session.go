package configs

import (
	"github.com/lemmego/api/config"
	"net/http"
)

var session = config.M{
	// Supported: "file", "database", "redis"
	"driver": config.MustEnv("SESSION_DRIVER", "file"),

	// Applicable when the driver is set to "database" or "redis"
	"connection": config.MustEnv("SESSION_CONNECTION", ""),

	"cookie": config.MustEnv("SESSION_COOKIE", "lemmego") + "_session",

	// Applicable when the driver is set to "file"
	"files": "./storage/session",

	"http_only": config.MustEnv("SESSION_HTTP_ONLY", true),
	"secure":    config.MustEnv("SESSION_SECURE_COOKIE", false),
	"domain":    config.MustEnv("SESSION_DOMAIN", ""),
	"path":      config.MustEnv("SESSION_PATH", "/"),
	"same_site": config.MustEnv("SESSION_SAME_SITE", http.SameSiteLaxMode),
}
