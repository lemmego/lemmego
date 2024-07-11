package handlers

import (
	"fmt"
	"lemmego/api"
)

// func IndexHandler(i *inertia.Inertia) http.Handler {

// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		err := i.Render(w, r, "Home/Index", nil)
// 		if err != nil {
// 			handleServerErr(w, err)
// 			return
// 		}
// 	}

// 	return http.HandlerFunc(fn)
// }

// func WelcomeHandler(i *inertia.Inertia) http.Handler {

// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		err := i.Render(w, r, "Home/Welcome", map[string]any{
// 			"name": "Tanmay",
// 		})
// 		if err != nil {
// 			handleServerErr(w, err)
// 			return
// 		}
// 	}

// 	return http.HandlerFunc(fn)
// }

// func handleServerErr(w http.ResponseWriter, err error) {
// 	log.Printf("http error: %s\n", err)
// 	w.WriteHeader(http.StatusInternalServerError)
// 	w.Write([]byte("server error"))
// }

func Register(r *api.Router) {
	//app.HTTPRouter().Get("/", app.Inertia().Middleware(IndexHandler(app.Inertia())).ServeHTTP)
	//app.HTTPRouter().Get("/welcome", app.Inertia().Middleware(WelcomeHandler(app.i)).ServeHTTP)
	//app.Post("/test", StoreTestHandler)
	r.Get("/", IndexHomeHandler, func(next api.Handler) api.Handler {
		return func(c *api.Context) error {
			return next(c)
		}
	})
	r.Group("/api", func(r *api.Router) {
		r.Group("/v1", func(r *api.Router) {
			r.Get("/test", func(c *api.Context) error {
				fmt.Println("Malta hit")
				return c.JSON(200, api.M{
					"foo": "bar",
				})
			})
			r.Get("/test2", func(c *api.Context) error {
				fmt.Println("Malta hit 2")
				return c.JSON(200, api.M{
					"foo": "bar",
				})
			})
		})
	}, func(next api.Handler) api.Handler {
		println("out")
		return func(c *api.Context) error {
			println("in")
			return next(c)
		}
	})
	//app.Get("/test2", func(c *api.Context) error {
	//	fmt.Println("Malta hit 2")
	//	return c.JSON(200, api.M{
	//		"foo": "bar",
	//	})
	//})

	//
	//app.Get("/login", LoginIndexHandler)
	//app.Post("/login", LoginStoreHandler)
	//app.Get("/register", RegistrationIndexHandler)
	//app.Post("/register", RegistrationStoreHandler)
}
