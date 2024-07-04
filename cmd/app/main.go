package main

import (
	"log/slog"
	"pressebo/api"
	_ "pressebo/internal/config"
	"pressebo/internal/handlers"
	"pressebo/internal/plugins"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Print config
	slog.Info("App will start using the following config:\n", "config", api.ConfMap())

	// Load plugins
	registry := plugins.Load()

	// Create application
	app := api.NewApp(
		api.WithPlugins(registry),
	)

	// Register global middleware
	app.Use(middleware.Logger, middleware.Recoverer)

	// Register routes
	handlers.Register(app)

	// Handle signals
	go app.HandleSignals()

	// Run application
	app.Run()
}
