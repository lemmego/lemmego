package config

import (
	"github.com/lemmego/lemmego/api"
	"github.com/lemmego/lemmego/api/config"
)

func init() {
	// Database
	config.SetConfig("db.driver", api.MustEnv("DB_DRIVER", "mysql"))
	config.SetConfig("db.host", api.MustEnv("DB_HOST", "localhost"))
	config.SetConfig("db.port", api.MustEnv("DB_PORT", 3306))
	config.SetConfig("db.database", api.MustEnv("DB_DATABASE", "pressebo"))
	config.SetConfig("db.username", api.MustEnv("DB_USERNAME", "root"))
	config.SetConfig("db.password", api.MustEnv("DB_PASSWORD", ""))
	config.SetConfig("db.params", api.MustEnv("DB_PARAMS", ""))

	// Redis
	config.SetConfig("db.redisHost", api.MustEnv("REDIS_HOST", "localhost"))
	config.SetConfig("db.redisPort", api.MustEnv("REDIS_PORT", 6379))
}
