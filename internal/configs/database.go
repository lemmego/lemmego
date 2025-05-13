package configs

import (
	"github.com/lemmego/api/config"
	"time"
)

func init() {
	config.Set("database", config.M{
		"default": config.MustEnv("DB_CONNECTION", "sqlite"),
		"connections": config.M{
			"sqlite": config.M{
				"driver":                  "sqlite",
				"url":                     config.MustEnv("DATABASE_URL", "file:./storage/database.sqlite?cache=shared&mode=memory"),
				"database":                config.MustEnv("DB_DATABASE", "./storage/database.sqlite"),
				"prefix":                  "",
				"foreign_key_constraints": config.MustEnv("DB_FOREIGN_KEYS", true),
			},
			"mysql": config.M{
				"driver":            "mysql",
				"host":              config.MustEnv("DB_HOST", "localhost"),
				"port":              config.MustEnv("DB_PORT", 3306),
				"database":          config.MustEnv("DB_DATABASE", "lemmego"),
				"user":              config.MustEnv("DB_USERNAME", "root"),
				"password":          config.MustEnv("DB_PASSWORD", ""),
				"params":            config.MustEnv("DB_PARAMS", ""),
				"auto_create":       config.MustEnv("DB_AUTOCREATE", false),
				"max_open_conns":    config.MustEnv("DB_MAX_OPEN_CONNS", 100),
				"max_idle_conns":    config.MustEnv("DB_MAX_IDLE_CONNS", 10),
				"conn_max_lifetime": config.MustEnv("DB_CONN_MAX_LIFETIME", time.Hour*2),
			},
			"pgsql": config.M{
				"driver":            "pgsql",
				"host":              config.MustEnv("DB_HOST", "localhost"),
				"port":              config.MustEnv("DB_PORT", 5432),
				"database":          config.MustEnv("DB_DATABASE", "lemmego"),
				"user":              config.MustEnv("DB_USERNAME", ""),
				"password":          config.MustEnv("DB_PASSWORD", ""),
				"params":            config.MustEnv("DB_PARAMS", ""),
				"auto_create":       config.MustEnv("DB_AUTOCREATE", false),
				"max_open_conns":    config.MustEnv("DB_MAX_OPEN_CONNS", 100),
				"max_idle_conns":    config.MustEnv("DB_MAX_IDLE_CONNS", 10),
				"conn_max_lifetime": config.MustEnv("DB_CONN_MAX_LIFETIME", time.Hour*2),
			},
		},
	})

	config.Set("redis", config.M{
		"connections": config.M{
			"default": config.M{
				"host":     config.MustEnv("REDIS_HOST", "localhost"),
				"port":     config.MustEnv("REDIS_PORT", 6379),
				"password": config.MustEnv("REDIS_PASSWORD", ""),
			},
		},
	})
}
