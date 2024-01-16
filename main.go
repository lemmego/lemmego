package main

import (
	"pressebo/framework"
	"pressebo/handlers"
	"pressebo/plugins"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	registry := plugins.Load()
	app := framework.NewApp(
		framework.WithPlugins(registry),
	)
	app.Use(middleware.Logger, middleware.Recoverer)

	handlers.RegisterRoutes(app)

	app.Run()
}
