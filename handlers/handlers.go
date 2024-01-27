package handlers

import (
	"pressebo/framework"
	// "pressebo/plugins/auth"
)

func RegisterRoutes(app *framework.App) {
	// auth := app.Plugin(auth.Namespace).(*auth.AuthPlugin)
	app.Get("/", IndexHomeHandler)
	app.Post("/test", StoreTestHandler)
}
