package middleware

import (
	"github.com/go-chi/httplog/v2"
	"github.com/lemmego/api/app"
	"log/slog"
)

func Logger() app.HTTPMiddleware {
	logger := httplog.NewLogger("lemmego", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		TimeFieldFormat:  "[15:04:05.000]",
	})

	return httplog.RequestLogger(logger)
}
