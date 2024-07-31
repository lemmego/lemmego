package providers

import (
	"fmt"
	"lemmego/api"
	"lemmego/internal/handlers"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

type RouteServiceProvider struct {
	*api.BaseServiceProvider
}

func (provider *RouteServiceProvider) Register(app *api.App) {
	logger := httplog.NewLogger("lemmego", httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		TimeFieldFormat:  "[15:04:05.000]",
		// Tags: map[string]string{
		// 	"version": "v1.0-81aa4244d9fc8076a",
		// 	"env":     "dev",
		// },
		// QuietDownRoutes: []string{
		// 	"/",
		// 	"/ping",
		// },
		// QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})

	// Routes global middleware
	app.Router().Use(httplog.RequestLogger(logger), middleware.Recoverer, func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
			return
		})
	})

	app.Router().UseBefore(func(next api.Handler) api.Handler {
		return func(c *api.Context) error {
			c.Set("foo", "bar")
			fmt.Println("I execute for every route")
			return next(c)
		}
	})

	app.RegisterRoutes(func(r *api.Router) {
		handlers.Routes(r)

		r.Get("/error", func(c *api.Context) error {
			err := c.SessionPop("error").(string)
			return c.HTML(500, "<html><body><code>"+err+"</code></body></html>")
		})
	})
}

func (provider *RouteServiceProvider) Boot() {
	//
}
