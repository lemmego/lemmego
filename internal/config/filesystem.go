package config

import (
	"github.com/lemmego/api/config"
)

func init() {
	config.Set("storage", map[string]any{
		"local": map[string]any{
			"driver": "local",
			"root":   "storage",
		},
		"s3": map[string]any{
			"driver":   "s3",
			"key":      config.MustEnv("AWS_ACCESS_KEY_ID", ""),
			"secret":   config.MustEnv("AWS_SECRET_ACCESS_KEY", ""),
			"region":   config.MustEnv("AWS_DEFAULT_REGION", "us-east-1"),
			"bucket":   config.MustEnv("AWS_BUCKET", ""),
			"endpoint": config.MustEnv("AWS_ENDPOINT", ""),
		},
		"r2": map[string]any{
			"driver":   "s3",
			"key":      config.MustEnv("R2_ACCESS_KEY_ID", ""),
			"secret":   config.MustEnv("R2_SECRET_ACCESS_KEY", ""),
			"region":   config.MustEnv("R2_DEFAULT_REGION", "us-east-1"),
			"bucket":   config.MustEnv("R2_BUCKET", ""),
			"endpoint": config.MustEnv("R2_ENDPOINT", ""),
		},
	})
}
