package main

import (
	"github.com/lemmego/api/app"
	_ "github.com/lemmego/api/logger"
	_ "github.com/lemmego/api/providers"
	//_ "github.com/lemmego/auth"
	"github.com/lemmego/lemmego/internal/commands"
	"github.com/lemmego/lemmego/internal/configs"
	_ "github.com/lemmego/lemmego/internal/migrations"
	"github.com/lemmego/lemmego/internal/routes"
)

func main() {
	// Configure application
	webApp := app.Configure(
		app.WithConfig(configs.Load()),
		app.WithCommands(commands.Load()),
		app.WithRoutes(routes.Load()),
	)

	// Run application
	webApp.Run()
}
