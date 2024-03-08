package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"pressebo/api"
	_ "pressebo/internal/config"
	"pressebo/internal/handlers"
	"pressebo/internal/plugins"
)

func main() {
	registry := plugins.Load()
	app := api.NewApp(
		api.WithPlugins(registry),
	)
	app.Use(middleware.Logger, middleware.Recoverer)

	handlers.RegisterRoutes(app)

	// Handle signals
	go app.HandleSignals()

	// Run application
	app.Run()
}
