package main

import (
	"fluent-blog/fluent"
	"fluent-blog/plugins/auth"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	authPlugin := auth.New(auth.WithResolveUser(func(username string, password string) (*auth.AuthUser, error) {
		// var db fluent.DBSession
		authUser := &auth.AuthUser{}
		// q := db.SQL().Select("id", "email", "password").From("users").Where("email = ?", username).Limit(1)
		// if err := q.One(&authUser); err != nil {
		// 	return nil, err
		// }
		// return authUser, nil
		authUser.ID = "1"
		authUser.Username = "tanmaymishu@gmail.com"
		authUser.Password = "$2y$10$koOG0SuI4WKbCgrjCocFPOg6OK7JH.Md6kPJX5EimoAQIOO8Nrs2m"
		return authUser, nil
	}))
	app := fluent.NewApp(fluent.WithPlugins([]fluent.Plugin{authPlugin}))
	app.Router.Use(middleware.Logger, middleware.Recoverer)

	welcomeHandler := func(ctx *fluent.Context) error {
		return ctx.Render(200, "welcome.page.tmpl", &fluent.TemplateData{})
	}

	app.Get("/", authPlugin.Web(welcomeHandler))

	app.Run()
}
