package handlers

import (
	"log"
	"net/http"
	"pressebo/api"

	inertia "github.com/romsar/gonertia"
)

func IndexHandler(i *inertia.Inertia) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		err := i.Render(w, r, "Home/Index", nil)
		if err != nil {
			handleServerErr(w, err)
			return
		}
	}

	return http.HandlerFunc(fn)
}

func WelcomeHandler(i *inertia.Inertia) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		err := i.Render(w, r, "Home/Welcome", map[string]any{
			"name": "Tanmay",
		})
		if err != nil {
			handleServerErr(w, err)
			return
		}
	}

	return http.HandlerFunc(fn)
}

func handleServerErr(w http.ResponseWriter, err error) {
	log.Printf("http error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("server error"))
}

func Register(app *api.App) {
	app.Router().Get("/", app.I.Middleware(IndexHandler(app.I)).ServeHTTP)
	app.Router().Get("/welcome", app.I.Middleware(WelcomeHandler(app.I)).ServeHTTP)
	// app.Get("/", IndexHomeHandler)
	app.Post("/test", StoreTestHandler)

	app.Get("/login", LoginIndexHandler)
	app.Post("/login", LoginStoreHandler)
	app.Get("/register", RegistrationIndexHandler)
	app.Post("/register", RegistrationStoreHandler)
}
