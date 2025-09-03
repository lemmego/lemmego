package bootstrap

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/middleware"
)

func LoadMiddlewares() []app.Handler {
	return []app.Handler{
		middleware.VerifyCSRF,
	}
}
