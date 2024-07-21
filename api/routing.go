package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/core"
	"github.com/go-chi/chi/v5"
	"github.com/golobby/container/v3"
	inertia "github.com/romsar/gonertia"
	"lemmego/api/logger"
	"lemmego/api/vee"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

const InKey = "input"

type Handler func(c *Context) error

type Middleware func(next Handler) Handler

func NewChiRouter() chi.Router {
	return chi.NewRouter()
}

type RouteRegistrarFunc func(r *Router)

type Router struct {
	chi.Router
	routeRegistrar   RouteRegistrarFunc
	currentGroup     *RouteGroup
	routeMiddlewares map[string]Middleware
	httpMiddlewares  []Middleware
	routes           []*Route
	basePrefix       string
}

// NewRouter creates a new HTTPRouter-based router
func NewRouter(router chi.Router) *Router {
	return &Router{router, nil, nil, nil, nil, nil, ""}
}

func (r *Router) setRouteMiddleware(middlewares map[string]Middleware) {
	r.routeMiddlewares = middlewares
}

func (hr *Router) Get(pattern string, handler Handler) *Route {
	fullPath := pattern
	if hr.currentGroup != nil {
		fullPath = path.Join(hr.currentGroup.prefix, pattern)
	}
	route := &Route{Method: http.MethodGet, Path: fullPath, Handler: handler}
	if hr.currentGroup != nil {
		hr.currentGroup.routes = append(hr.currentGroup.routes, route)
	}
	hr.routes = append(hr.routes, route)
	return route
}

func (hr *Router) Post(pattern string, handler Handler) *Route {
	fullPath := pattern
	if hr.currentGroup != nil {
		fullPath = path.Join(hr.currentGroup.prefix, pattern)
	}
	route := &Route{Method: http.MethodPost, Path: fullPath, Handler: handler}
	if hr.currentGroup != nil {
		hr.currentGroup.routes = append(hr.currentGroup.routes, route)
	}
	hr.routes = append(hr.routes, route)
	return route
}

func (hr *Router) Put(pattern string, handler Handler) *Route {
	fullPath := pattern
	if hr.currentGroup != nil {
		fullPath = path.Join(hr.currentGroup.prefix, pattern)
	}
	route := &Route{Method: http.MethodPut, Path: fullPath, Handler: handler}
	if hr.currentGroup != nil {
		hr.currentGroup.routes = append(hr.currentGroup.routes, route)
	}
	hr.routes = append(hr.routes, route)
	return route
}

func (hr *Router) Patch(pattern string, handler Handler) *Route {
	fullPath := pattern
	if hr.currentGroup != nil {
		fullPath = path.Join(hr.currentGroup.prefix, pattern)
	}
	route := &Route{Method: http.MethodPatch, Path: fullPath, Handler: handler}
	if hr.currentGroup != nil {
		hr.currentGroup.routes = append(hr.currentGroup.routes, route)
	}
	hr.routes = append(hr.routes, route)
	return route
}

func (hr *Router) Delete(pattern string, handler Handler) *Route {
	fullPath := pattern
	if hr.currentGroup != nil {
		fullPath = path.Join(hr.currentGroup.prefix, pattern)
	}
	route := &Route{Method: http.MethodDelete, Path: fullPath, Handler: handler}
	if hr.currentGroup != nil {
		hr.currentGroup.routes = append(hr.currentGroup.routes, route)
	}
	hr.routes = append(hr.routes, route)
	return route
}

func (hr *Router) Connect(pattern string, handler Handler) *Route {
	fullPath := pattern
	if hr.currentGroup != nil {
		fullPath = path.Join(hr.currentGroup.prefix, pattern)
	}
	route := &Route{Method: http.MethodConnect, Path: fullPath, Handler: handler}
	if hr.currentGroup != nil {
		hr.currentGroup.routes = append(hr.currentGroup.routes, route)
	}
	hr.routes = append(hr.routes, route)
	return route
}

func (hr *Router) Head(pattern string, handler Handler) *Route {
	fullPath := pattern
	if hr.currentGroup != nil {
		fullPath = path.Join(hr.currentGroup.prefix, pattern)
	}
	route := &Route{Method: http.MethodHead, Path: fullPath, Handler: handler}
	if hr.currentGroup != nil {
		hr.currentGroup.routes = append(hr.currentGroup.routes, route)
	}
	hr.routes = append(hr.routes, route)
	return route
}

func (hr *Router) Options(pattern string, handler Handler) *Route {
	fullPath := pattern
	if hr.currentGroup != nil {
		fullPath = path.Join(hr.currentGroup.prefix, pattern)
	}
	route := &Route{Method: http.MethodOptions, Path: fullPath, Handler: handler}
	if hr.currentGroup != nil {
		hr.currentGroup.routes = append(hr.currentGroup.routes, route)
	}
	hr.routes = append(hr.routes, route)
	return route
}

func (hr *Router) Trace(pattern string, handler Handler) *Route {
	fullPath := pattern
	if hr.currentGroup != nil {
		fullPath = path.Join(hr.currentGroup.prefix, pattern)
	}
	route := &Route{Method: http.MethodTrace, Path: fullPath, Handler: handler}
	if hr.currentGroup != nil {
		hr.currentGroup.routes = append(hr.currentGroup.routes, route)
	}
	hr.routes = append(hr.routes, route)
	return route
}

// RouteGroup represents a group of routes
type RouteGroup struct {
	prefix           string
	beforeMiddleware []Middleware
	afterMiddleware  []Middleware
	routes           []*Route
}

func (rg *RouteGroup) UseBefore(middlewares ...Middleware) *RouteGroup {
	rg.beforeMiddleware = append(rg.beforeMiddleware, middlewares...)
	for _, route := range rg.routes {
		route.BeforeMiddleware = append(middlewares, route.BeforeMiddleware...)
	}
	return rg
}

func (rg *RouteGroup) UseAfter(middlewares ...Middleware) *RouteGroup {
	rg.afterMiddleware = append(rg.afterMiddleware, middlewares...)
	for _, route := range rg.routes {
		route.AfterMiddleware = append(route.AfterMiddleware, middlewares...)
	}
	return rg
}

// Group method for App
func (hr *Router) Group(prefix string, fn func(r *Router)) *RouteGroup {
	previousGroup := hr.currentGroup
	newGroup := &RouteGroup{
		prefix: prefix,
	}

	if previousGroup != nil {
		newGroup.prefix = path.Join(previousGroup.prefix, newGroup.prefix)
	}

	hr.currentGroup = newGroup

	subRouter := &Router{
		Router:         hr.Router,
		routeRegistrar: hr.routeRegistrar,
		currentGroup:   newGroup,
	}

	fn(subRouter)

	newGroup.routes = subRouter.routes
	hr.routes = append(hr.routes, subRouter.routes...)

	hr.currentGroup = previousGroup

	return newGroup
}

func (hr *Router) registerMiddlewares(app *App) {
	for _, plugin := range app.plugins {
		hr.httpMiddlewares = append(hr.httpMiddlewares, plugin.Middlewares()...)
	}
	for _, mw := range hr.httpMiddlewares {
		hr.Use(adaptMiddleware(app, mw))
	}
}

func (hr *Router) registerRoutes(app *App) {
	for _, plugin := range app.plugins {
		for _, route := range plugin.Routes() {
			hr.routes = append(hr.routes, route)
		}
	}
	hr.routeRegistrar(hr)

	for _, route := range hr.routes {
		handler := route.Handler

		// Apply group middlewares
		if hr.currentGroup != nil {
			for i := len(hr.currentGroup.beforeMiddleware) - 1; i >= 0; i-- {
				handler = hr.currentGroup.beforeMiddleware[i](handler)
			}
			for i := len(hr.currentGroup.afterMiddleware) - 1; i >= 0; i-- {
				nextHandler := handler
				handler = func(c *Context) error {
					err := nextHandler(c)
					if err != nil {
						return err
					}
					return hr.currentGroup.afterMiddleware[i](func(*Context) error { return nil })(c)
				}
			}
		}

		hr.MethodFunc(route.Method, route.Path, makeHandlerFunc(app, &Route{
			Method:           route.Method,
			Path:             route.Path,
			Handler:          handler,
			BeforeMiddleware: route.BeforeMiddleware,
			AfterMiddleware:  route.AfterMiddleware,
		}))
	}
}

type Route struct {
	Method           string
	Path             string
	Handler          Handler
	AfterMiddleware  []Middleware
	BeforeMiddleware []Middleware
}

func (r *Route) UseBefore(middleware ...Middleware) *Route {
	r.BeforeMiddleware = append(r.BeforeMiddleware, middleware...)
	return r
}

func (r *Route) UseAfter(middleware ...Middleware) *Route {
	r.AfterMiddleware = append(r.AfterMiddleware, middleware...)
	return r
}

func adaptMiddleware(app *App, m Middleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a custom Handler that wraps the next http.Handler
			customNext := Handler(func(c *Context) error {
				next.ServeHTTP(w, r)
				return nil
			})

			// Apply the custom middleware
			wrappedHandler := m(customNext)

			// Create a new Context (you'll need to adjust this based on your Context structure)
			ctx := &Context{sync.Mutex{}, app, container.New(), r, w}

			// Call the wrapped handler
			if err := wrappedHandler(ctx); err != nil {
				// Handle error (you might want to adjust this based on your error handling strategy)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
}

func makeHandlerFunc(app *App, route *Route) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{sync.Mutex{}, app, container.New(), r, w}
		if !app.isContextReady {
			app.isContextReady = true
		}

		chain := route.Handler

		// Apply after middlewares in reverse order
		for i := len(route.AfterMiddleware) - 1; i >= 0; i-- {
			afterMiddleware := route.AfterMiddleware[i]
			nextChain := chain
			chain = func(c *Context) error {
				err := nextChain(c)
				if err != nil {
					return err
				}
				return afterMiddleware(func(*Context) error { return nil })(c)
			}
		}

		// Apply before middlewares
		for i := len(route.BeforeMiddleware) - 1; i >= 0; i-- {
			chain = route.BeforeMiddleware[i](chain)
		}

		if err := chain(ctx); err != nil {
			logger.V().Error(err.Error())
			// Check if the error is validation errors, if so
			// return 422 according to the UI preference.
			if errors.As(err, &vee.Errors{}) {
				if !ctx.WantsJSON() {
					ctx.WithErrors(err.(vee.Errors)).Back(http.StatusFound)
					return
				}
				err := ctx.JSON(http.StatusUnprocessableEntity, M{"errors": err})
				if err != nil {
					ctx.Error(http.StatusInternalServerError, err)
				}
				return
			}

			if !ctx.WantsJSON() {
				ctx.WithError(err.Error()).
					Redirect(http.StatusFound, "/error")
				return
			}
			ctx.Error(http.StatusInternalServerError, err)
			return
		}
	}

	if app.i != nil {
		return app.Inertia().Middleware(http.HandlerFunc(fn)).ServeHTTP
	}

	return fn
}

func initInertia() *inertia.Inertia {
	manifestPath := "./public/build/manifest.json"

	i, err := inertia.NewFromFile(
		"resources/views/root.html",
		// inertia.WithVersion("1.0"),
		inertia.WithVersionFromFile(manifestPath),
		inertia.WithSSR(),
	)

	if err != nil {
		log.Fatal(err)
	}

	i.ShareTemplateFunc("vite", vite(manifestPath, "/public/build/"))
	i.ShareTemplateData("env", Config("app.env").(string))

	return i
}

func vite(manifestPath, buildDir string) func(path string) (string, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		log.Fatalf("cannot open provided vite manifest file: %s", err)
	}
	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
	})
	err = json.NewDecoder(f).Decode(&viteAssets)
	if err != nil {
		log.Fatalf("cannot unmarshal vite manifest file to json: %s", err)
	}

	return func(p string) (string, error) {
		if val, ok := viteAssets[p]; ok {
			return path.Join(buildDir, val.File), nil
		}
		return "", fmt.Errorf("asset %q not found", p)
	}
}

//func NewRoute(method string, path string, handler Handler, middlewares ...Middleware) *Route {
//	return &Route{
//		Method:      method,
//		Path:        path,
//		Middlewares: middlewares,
//		Handler:     handler,
//	}
//}

//	func Input(inputStruct any, opts ...core.Option) Middleware {
//		co, err := httpin.New(inputStruct, opts...)
//
//		if err != nil {
//			panic(err)
//		}
//
//		return func(next http.Handler) http.Handler {
//			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				// Create a context (you might need to adjust this based on your actual Context creation logic)
//				ctx := &Context{
//					responseWriter: w,
//					request:        r,
//					// Initialize other fields as necessary
//				}
//
//				input, err := co.Decode(r)
//				if err != nil {
//					co.GetErrorHandler()(w, r, err)
//					return
//				}
//
//				ctx.Set(InKey, input)
//
//				// Create a new request with the context
//				r = r.WithContext(context.WithValue(r.Context(), InKey, ctx))
//
//				// Call the next handler
//				next.ServeHTTP(w, r)
//			})
//		}
//	}

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
