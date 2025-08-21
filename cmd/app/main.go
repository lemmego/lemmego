package main

import (
	"github.com/lemmego/api/app"
	_ "github.com/lemmego/api/logger"
	_ "github.com/lemmego/api/providers"
	"github.com/lemmego/lemmego/bootstrap"
	_ "github.com/lemmego/lemmego/internal/commands"
	_ "github.com/lemmego/lemmego/internal/configs"
	_ "github.com/lemmego/lemmego/internal/middleware"
	_ "github.com/lemmego/lemmego/internal/migrations"
	_ "github.com/lemmego/lemmego/internal/routes"
)

func main() {
	// Configure an instance of the application
	webApp := app.Configure()

	// Bootstrap application
	webApp.WithRoutes(bootstrap.LoadRoutes()).
		WithHTTPMiddlewares(bootstrap.LoadHTTPMiddlewares()).
		WithMiddlewares(bootstrap.LoadMiddlewares()).
		WithCommands(bootstrap.LoadCommands()).
		WithProviders(bootstrap.LoadProviders())

	// Run application
	webApp.Run()
}
