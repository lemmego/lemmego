package main

import (
	"log/slog"

	"github.com/lemmego/lemmego/api/app"
	"github.com/lemmego/lemmego/api/config"
	_ "github.com/lemmego/lemmego/internal/config"
	"github.com/lemmego/lemmego/internal/plugins"
	"github.com/lemmego/lemmego/internal/providers"
)

func main() {
	// Print config
	slog.Info("App will start using the following config:\n", "config", config.ConfMap())

	// Load service providers
	providerCollection := providers.Load()
	// Load plugins
	pluginCollection := plugins.Load()

	// Create application
	app := app.NewApp(
		app.WithProviders(providerCollection),
		app.WithPlugins(pluginCollection),
		app.WithInertia(nil),
		app.WithFS(nil),
	)

	// Handle signals
	go app.HandleSignals()

	// Run application
	app.Run()
}
