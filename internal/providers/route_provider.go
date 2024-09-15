package providers

import (
	"log/slog"

	"github.com/lemmego/api/app"
	"github.com/lemmego/api/middleware"
	"github.com/lemmego/lemmego/internal/handlers"

	chiMw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

type RouteServiceProvider struct {
	*app.BaseServiceProvider
}

func (provider *RouteServiceProvider) Register(a *app.App) {
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

	// net/http compatible global middleware
	a.Router().Use(httplog.RequestLogger(logger), chiMw.Recoverer)

	// Global middleware
	a.Router().UseBefore(middleware.VerifyCSRF)

	a.RegisterRoutes(func(r *app.Router) {
		handlers.Routes(r)

		r.Get("/error", func(c *app.Context) error {
			err := c.PopSession("error").(string)
			return c.HTML(500, []byte("<html><body><code>"+err+"</code></body></html>"))
		})
	})
}

func (provider *RouteServiceProvider) Boot() {
	//
}
