package config

import (
	"github.com/lemmego/lemmego/api/config"
	"github.com/lemmego/lemmego/api/utils"
)

func init() {
	// Database
	config.SetConfig("db.driver", utils.MustEnv("DB_DRIVER", "mysql"))
	config.SetConfig("db.host", utils.MustEnv("DB_HOST", "localhost"))
	config.SetConfig("db.port", utils.MustEnv("DB_PORT", 3306))
	config.SetConfig("db.database", utils.MustEnv("DB_DATABASE", "pressebo"))
	config.SetConfig("db.username", utils.MustEnv("DB_USERNAME", "root"))
	config.SetConfig("db.password", utils.MustEnv("DB_PASSWORD", ""))
	config.SetConfig("db.params", utils.MustEnv("DB_PARAMS", ""))

	// Redis
	config.SetConfig("db.redisHost", utils.MustEnv("REDIS_HOST", "localhost"))
	config.SetConfig("db.redisPort", utils.MustEnv("REDIS_PORT", 6379))
}
