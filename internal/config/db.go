package config

import (
	"pressebo/api"
)

func init() {
	api.SetConfig("db.driver", api.MustEnv("DB_DRIVER", "mysql"))
	api.SetConfig("db.host", api.MustEnv("DB_HOST", "localhost"))
	api.SetConfig("db.port", api.MustEnv("DB_PORT", 3306))
	api.SetConfig("db.database", api.MustEnv("DB_DATABASE", "pressebo"))
	api.SetConfig("db.username", api.MustEnv("DB_USERNAME", "root"))
	api.SetConfig("db.password", api.MustEnv("DB_PASSWORD", ""))
	api.SetConfig("db.params", api.MustEnv("DB_PARAMS", ""))
}
