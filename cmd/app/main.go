package main

import (
	"lemmego/api"
	_ "lemmego/internal/config"
	"lemmego/internal/plugins"
	"lemmego/internal/providers"
	"log/slog"
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
