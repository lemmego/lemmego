package handlers

import (
	"pressebo/api"
)

func Register(app *api.App) {
	app.Get("/", IndexHomeHandler)
	app.Post("/test", StoreTestHandler)

	app.Get("/login", LoginIndexHandler)
	app.Post("/login", LoginStoreHandler)
	app.Get("/register", RegistrationIndexHandler)
	app.Post("/register", RegistrationStoreHandler)
}
