package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"lemmego/api"
	_ "lemmego/internal/config"
	"lemmego/internal/handlers"
	"lemmego/internal/plugins"
	"log/slog"
)

func main() {
	// Print config
	slog.Info("App will start using the following config:\n", "config", api.ConfMap())

	// Load plugins
	registry := plugins.Load()

	// Create application
	app := api.NewApp(
		api.WithPlugins(registry),
		api.WithInertia(nil),
		api.WithFS(nil),
	)

	logger := httplog.NewLogger("lemmego", httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		TimeFieldFormat:  "[15:04:05.000]",
		// Tags: map[string]string{
		// 	"version": "v1.0-81aa4244d9fc8076a",
		// 	"env":     "dev",
		// },
		// QuietDownRoutes: []string{
		// 	"/",
		// 	"/ping",
		// },
		// QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})

	// Register global middleware
	app.Router().Use(httplog.RequestLogger(logger), middleware.Recoverer)

	app.RegisterRoutes(func(r *api.Router) {
		handlers.Register(r)
	})

	// Handle signals
	go app.HandleSignals()

	// Run application
	app.Run()
}
