package config

import (
	"github.com/lemmego/lemmego/api/config"
)

func init() {
	config.Set("database", map[string]any{
		"connections": map[string]any{
			"default": map[string]any{
				"driver":   config.MustEnv("DB_DRIVER", "mysql"),
				"host":     config.MustEnv("DB_HOST", "localhost"),
				"port":     config.MustEnv("DB_PORT", 3306),
				"database": config.MustEnv("DB_DATABASE", "pressebo"),
				"user":     config.MustEnv("DB_USERNAME", "root"),
				"password": config.MustEnv("DB_PASSWORD", ""),
				"params":   config.MustEnv("DB_PARAMS", ""),
			},
		},
	})

	config.Set("redis", map[string]any{
		"connections": map[string]any{
			"default": map[string]any{
				"host":     config.MustEnv("REDIS_HOST", "localhost"),
				"port":     config.MustEnv("REDIS_PORT", 6379),
				"password": config.MustEnv("REDIS_PASSWORD", ""),
			},
		},
	})
}
