package configs

import (
	"github.com/lemmego/api/config"
)

func Load() config.M {
	return config.M{
		"app":         app,
		"database":    database["database"],
		"redis":       database["redis"],
		"filesystems": filesystems,
	}
}
