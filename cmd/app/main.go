package main

import (
	"github.com/lemmego/api/app"
	_ "github.com/lemmego/api/logger"
	"github.com/lemmego/lemmego/internal/commands"
	"github.com/lemmego/lemmego/internal/configs"
	_ "github.com/lemmego/lemmego/internal/migrations"
	"github.com/lemmego/lemmego/internal/plugins"
	"github.com/lemmego/lemmego/internal/providers"
	"github.com/lemmego/lemmego/internal/routes"
)

func main() {
	// Print config
	//slog.Info("app will start using the following config:\n", "config", config.GetAll())
	// Create application
	webApp := app.New()

	webApp.WithConfig(configs.Load()).
		WithCommands(commands.Load()).
		WithProviders(providers.Load()).
		WithPlugins(plugins.Load()).
		WithRoutes(routes.Load())

	// Handle signals
	go webApp.HandleSignals()

	// Run application
	webApp.Run()
}
