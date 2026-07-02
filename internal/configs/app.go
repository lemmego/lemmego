package configs

import "github.com/lemmego/api/config"

func init() {
	config.Set("app", config.M{
		"name":            config.MustEnv("APP_NAME", "Lemmego"),
		"port":            config.MustEnv("APP_PORT", 8080),
		"env":             config.MustEnv("APP_ENV", "development"),
		"debug":           config.MustEnv("APP_DEBUG", false),
		"config_path":     config.MustEnv("CONFIG_PATH", "./internal/configs"),
		"command_path":    config.MustEnv("COMMAND_PATH", "./internal/commands"),
		"handler_path":    config.MustEnv("HANDLER_PATH", "./internal/handlers"),
		"input_path":      config.MustEnv("INPUT_PATH", "./internal/inputs"),
		"middleware_path": config.MustEnv("MIDDLEWARE_PATH", "./internal/middlewares"),
		"migration_path":  config.MustEnv("MIGRATION_PATH", "./internal/migrations"),
		"model_path":      config.MustEnv("MODEL_PATH", "./internal/models"),
		"route_path":      config.MustEnv("ROUTE_PATH", "./internal/routes"),
	})
}
