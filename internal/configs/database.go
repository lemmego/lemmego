package configs

import "github.com/lemmego/api/config"

var database = config.M{
	"database": config.M{
		"default": config.MustEnv("DB_CONNECTION", "sqlite"),
		"connections": config.M{
			"mysql": config.M{
				"driver":   "mysql",
				"host":     config.MustEnv("DB_HOST", "localhost"),
				"port":     config.MustEnv("DB_PORT", 3306),
				"database": config.MustEnv("DB_DATABASE", "lemmego"),
				"user":     config.MustEnv("DB_USERNAME", "root"),
				"password": config.MustEnv("DB_PASSWORD", ""),
				"params":   config.MustEnv("DB_PARAMS", ""),
			},
			"pgsql": config.M{
				"driver":   "pgsql",
				"host":     config.MustEnv("DB_HOST", "localhost"),
				"port":     config.MustEnv("DB_PORT", 5432),
				"database": config.MustEnv("DB_DATABASE", "lemmego"),
				"user":     config.MustEnv("DB_USERNAME", ""),
				"password": config.MustEnv("DB_PASSWORD", ""),
				"params":   config.MustEnv("DB_PARAMS", ""),
			},
		},
	},
	"redis": config.M{
		"connections": config.M{
			"default": config.M{
				"host":     config.MustEnv("REDIS_HOST", "localhost"),
				"port":     config.MustEnv("REDIS_PORT", 6379),
				"password": config.MustEnv("REDIS_PASSWORD", ""),
			},
		},
	},
}
