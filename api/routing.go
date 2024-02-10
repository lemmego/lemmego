package api

import (
	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/core"
	"github.com/go-chi/chi/v5"
)

const InKey = "input"

type Router interface {
	chi.Router
}

type router struct {
	chi.Router
}

func NewRouter() Router {
	return &router{chi.NewRouter()}
}

type Route struct {
	Method      string
	Path        string
	Middlewares []Middleware
	Handler     Handler
}

func NewRoute(method string, path string, handler Handler, middlewares ...Middleware) *Route {
	return &Route{
		Method:      method,
		Path:        path,
		Middlewares: middlewares,
		Handler:     handler,
	}
}

func Input(inputStruct any, opts ...core.Option) Middleware {
	co, err := httpin.New(inputStruct, opts...)

	if err != nil {
		panic(err)
	}

	return func(next Handler) Handler {
		return func(ctx *Context) error {
			input, err := co.Decode(ctx.Request())
			if err != nil {
				co.GetErrorHandler()(ctx.ResponseWriter(), ctx.Request(), err)
				return nil
			}

			ctx.Set(InKey, input)
			return next(ctx)
		}
	}
}
