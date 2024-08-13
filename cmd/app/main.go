package main

import (
	"log/slog"

	"github.com/lemmego/lemmego/api"
	_ "github.com/lemmego/lemmego/internal/config"
	"github.com/lemmego/lemmego/internal/plugins"
	"github.com/lemmego/lemmego/internal/providers"
)

func main() {
	// Print config
	slog.Info("App will start using the following config:\n", "config", api.ConfMap())

	// Load service providers
	providerCollection := providers.Load()
	// Load plugins
	pluginCollection := plugins.Load()

	// Create application
	app := api.NewApp(
		api.WithProviders(providerCollection),
		api.WithPlugins(pluginCollection),
		api.WithInertia(nil),
		api.WithFS(nil),
	)

	// Handle signals
	go app.HandleSignals()

	// Run application
	app.Run()
}
