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

type HTTPRouter interface {
	http.Handler

	// Use appends one or more middlewares onto the HTTPRouter stack.
	Use(middlewares ...func(http.Handler) http.Handler)

	// With adds inline middlewares for an endpoint handler.
	With(middlewares ...func(http.Handler) http.Handler) HTTPRouter

	// Group adds a new inline-HTTPRouter along the current routing
	// path, with a fresh middleware stack for the inline-HTTPRouter.
	Group(fn func(r HTTPRouter)) HTTPRouter

	// Route mounts a sub-HTTPRouter along a `pattern`` string.
	Route(pattern string, fn func(r HTTPRouter)) HTTPRouter

	// Mount attaches another http.Handler along ./pattern/*
	Mount(pattern string, h http.Handler)

	// Handle and HandleFunc adds routes for `pattern` that matches
	// all HTTP methods.
	Handle(pattern string, h http.Handler)
	HandleFunc(pattern string, h http.HandlerFunc)

	// Method and MethodFunc adds routes for `pattern` that matches
	// the `method` HTTP method.
	Method(method, pattern string, h http.Handler)
	MethodFunc(method, pattern string, h http.HandlerFunc)

	// HTTP-method routing along `pattern`
	Connect(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
	Get(pattern string, h http.HandlerFunc)
	Head(pattern string, h http.HandlerFunc)
	Options(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Trace(pattern string, h http.HandlerFunc)

	// NotFound defines a handler to respond whenever a route could
	// not be found.
	NotFound(h http.HandlerFunc)

	// MethodNotAllowed defines a handler to respond whenever a method is
	// not allowed.
	MethodNotAllowed(h http.HandlerFunc)

	BaseRouter() interface{}
}

// Routes interface adds two methods for router traversal, which is also
// used by the `docgen` subpackage to generation documentation for Routers.
//type Routes interface {
//	// Routes returns the routing tree in an easily traversable structure.
//	Routes() []Route
//
//	// Middlewares returns the list of middlewares in use by the router.
//	Middlewares() Middlewares
//
//	// Match searches the routing tree for a handler that matches
//	// the method/path - similar to routing a http request, but without
//	// executing the handler thereafter.
//	Match(rctx *Context, method, path string) bool
//}

// Middlewares type is a slice of standard middleware handlers with methods
// to compose middleware chains and http.Handler's.
//type Middlewares []func(http.Handler) http.Handler

type router struct {
	HTTPRouter
}

type chiRouter struct {
	chi.Router
}

func NewChiRouter() HTTPRouter {
	return &chiRouter{chi.NewRouter()}
}

// Implement the methods of your HTTPRouter interface that need special handling
func (r *chiRouter) With(middlewares ...func(http.Handler) http.Handler) HTTPRouter {
	return &chiRouter{r.Router.With(middlewares...)}
}

func (r *chiRouter) Group(fn func(r HTTPRouter)) HTTPRouter {
	return &chiRouter{r.Router.Group(func(chiR chi.Router) {
		fn(&chiRouter{chiR})
	})}
}

func (r *chiRouter) Route(pattern string, fn func(r HTTPRouter)) HTTPRouter {
	return &chiRouter{r.Router.Route(pattern, func(chiR chi.Router) {
		fn(&chiRouter{chiR})
	})}
}

func (r *chiRouter) BaseRouter() interface{} {
	return r
}

type RouteRegistrarFunc func(r *Router)

type Router struct {
	HTTPRouter
	routeRegistrar   RouteRegistrarFunc
	currentGroup     *RouteGroup
	routeMiddlewares map[string]Middleware
	httpMiddlewares  []Middleware
	routes           []*Route
	basePrefix       string
}

// NewRouter creates a new HTTPRouter-based router
func NewRouter(router HTTPRouter) *Router {
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
		HTTPRouter:     hr.HTTPRouter,
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
	return func(w http.ResponseWriter, r *http.Request) {
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
			if !ctx.WantsJSON() {
				ctx.Redirect(http.StatusFound, "/error")
				return
			}
			if errors.As(err, &vee.Errors{}) {
				ctx.JSON(http.StatusInternalServerError, M{"errors": err})
			} else {
				ctx.JSON(http.StatusInternalServerError, M{"error": err.Error()})
			}
		}
	}
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
