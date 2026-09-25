package main

import (
	"log"
	"os"

	"github.com/lemmego/api/app"
	_ "github.com/lemmego/api/logger"
	"github.com/lemmego/api/res"
	"github.com/lemmego/lemmego/bootstrap"
	_ "github.com/lemmego/lemmego/internal/configs"
	_ "github.com/lemmego/lemmego/internal/migrations"
)

// templatesDir holds the page templates: the error pages, and any Go template
// views the project adds.
const templatesDir = "templates"

func main() {
	// The template cache is populated explicitly rather than on import, so
	// that importing the package never depends on the working directory.
	// Nothing did populate it, which left every .gohtml render failing with
	// "not found in cache" — the error pages included.
	if _, err := os.Stat(templatesDir); err == nil {
		if err := res.LoadTemplates(templatesDir); err != nil {
			log.Fatalf("loading templates from %s: %v", templatesDir, err)
		}
	}

	// Configure an instance of the application
	webApp := app.Configure()

	// Bootstrap application
	webApp.WithRoutes(bootstrap.LoadRoutes()).
		WithHTTPMiddlewares(bootstrap.LoadHTTPMiddlewares()).
		WithMiddlewares(bootstrap.LoadMiddlewares()).
		WithCommands(bootstrap.LoadCommands()).
		WithProviders(bootstrap.LoadProviders()).
		WithErrMap(bootstrap.LoadErrMap())

	// Run application
	webApp.Run()
}
