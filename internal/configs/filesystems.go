package configs

import (
	"github.com/lemmego/api/config"
)

var filesystems = config.M{
	"default": config.MustEnv("FILESYSTEM_DISK", "local"),
	"disks": config.M{
		"local": config.M{
			"driver": "local",
			"root":   "storage",
			"path":   "./storage",
		},
		"s3": config.M{
			"driver":   "s3",
			"key":      config.MustEnv("AWS_ACCESS_KEY_ID", ""),
			"secret":   config.MustEnv("AWS_SECRET_ACCESS_KEY", ""),
			"region":   config.MustEnv("AWS_DEFAULT_REGION", "us-east-1"),
			"bucket":   config.MustEnv("AWS_BUCKET", ""),
			"endpoint": config.MustEnv("AWS_ENDPOINT", ""),
		},
		"r2": config.M{
			"driver":   "s3",
			"key":      config.MustEnv("R2_ACCESS_KEY_ID", ""),
			"secret":   config.MustEnv("R2_SECRET_ACCESS_KEY", ""),
			"region":   config.MustEnv("R2_DEFAULT_REGION", "us-east-1"),
			"bucket":   config.MustEnv("R2_BUCKET", ""),
			"endpoint": config.MustEnv("R2_ENDPOINT", ""),
		},
	},
}
