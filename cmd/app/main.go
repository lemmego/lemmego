package main

import (
	"github.com/lemmego/api/app"
	_ "github.com/lemmego/api/logger"
	_ "github.com/lemmego/api/providers"
	_ "github.com/lemmego/lemmego/internal/commands"
	_ "github.com/lemmego/lemmego/internal/configs"
	_ "github.com/lemmego/lemmego/internal/middleware"
	_ "github.com/lemmego/lemmego/internal/migrations"
	_ "github.com/lemmego/lemmego/internal/routes"
)

func main() {
	// Configure application
	webApp := app.Configure()

	// Run application
	webApp.Run()
}
