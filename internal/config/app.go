package config

import (
	"log"
	"pressebo/api"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api.SetConfig("app.name", api.MustEnv("APP_NAME", "Pressebo"))
	api.SetConfig("app.port", api.MustEnv("APP_PORT", 8080))
	api.SetConfig("app.env", api.MustEnv("APP_ENV", "development"))
	api.SetConfig("app.debug", api.MustEnv("APP_DEBUG", "true"))
	api.SetConfig("app.templateDir", "templates")
}
