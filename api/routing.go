package api

import (
	"encoding/json"
	"fmt"
	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/core"
	"github.com/go-chi/chi/v5"
	"github.com/golobby/container/v3"
	inertia "github.com/romsar/gonertia"
	"lemmego/api/logger"
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
}

// NewRouter creates a new HTTPRouter-based router
func NewRouter(router HTTPRouter) *Router {
	return &Router{router, nil, nil, nil, nil, nil}
}

func (r *Router) setRouteMiddleware(middlewares map[string]Middleware) {
	r.routeMiddlewares = middlewares
}

func (hr *Router) Get(pattern string, handler Handler, middlewares ...Middleware) {
	fullPattern, fullMiddlewares := hr.appendPatternAndMiddlewares(pattern, middlewares)

	hr.routes = append(hr.routes, &Route{http.MethodGet, fullPattern, fullMiddlewares, handler})
}

func (hr *Router) Post(pattern string, handler Handler, middlewares ...Middleware) {
	fullPattern, fullMiddlewares := hr.appendPatternAndMiddlewares(pattern, middlewares)

	hr.routes = append(hr.routes, &Route{http.MethodPost, fullPattern, fullMiddlewares, handler})
}

func (hr *Router) Put(pattern string, handler Handler, middlewares ...Middleware) {
	fullPattern, fullMiddlewares := hr.appendPatternAndMiddlewares(pattern, middlewares)

	hr.routes = append(hr.routes, &Route{http.MethodPut, fullPattern, fullMiddlewares, handler})
}

func (hr *Router) Patch(pattern string, handler Handler, middlewares ...Middleware) {
	fullPattern, fullMiddlewares := hr.appendPatternAndMiddlewares(pattern, middlewares)

	hr.routes = append(hr.routes, &Route{http.MethodPatch, fullPattern, fullMiddlewares, handler})
}

func (hr *Router) Delete(pattern string, handler Handler, middlewares ...Middleware) {
	fullPattern, fullMiddlewares := hr.appendPatternAndMiddlewares(pattern, middlewares)

	hr.routes = append(hr.routes, &Route{http.MethodDelete, fullPattern, fullMiddlewares, handler})
}

func (hr *Router) Connect(pattern string, handler Handler, middlewares ...Middleware) {
	fullPattern, fullMiddlewares := hr.appendPatternAndMiddlewares(pattern, middlewares)

	hr.routes = append(hr.routes, &Route{http.MethodConnect, fullPattern, fullMiddlewares, handler})
}

func (hr *Router) Head(pattern string, handler Handler, middlewares ...Middleware) {
	fullPattern, fullMiddlewares := hr.appendPatternAndMiddlewares(pattern, middlewares)

	hr.routes = append(hr.routes, &Route{http.MethodHead, fullPattern, fullMiddlewares, handler})
}

func (hr *Router) Options(pattern string, handler Handler, middlewares ...Middleware) {
	fullPattern, fullMiddlewares := hr.appendPatternAndMiddlewares(pattern, middlewares)

	hr.routes = append(hr.routes, &Route{http.MethodOptions, fullPattern, fullMiddlewares, handler})
}

func (hr *Router) Trace(pattern string, handler Handler, middlewares ...Middleware) {
	fullPattern, fullMiddlewares := hr.appendPatternAndMiddlewares(pattern, middlewares)

	hr.routes = append(hr.routes, &Route{http.MethodTrace, fullPattern, fullMiddlewares, handler})
}

// Group method for App
func (hr *Router) Group(prefix string, fn func(r *Router), middlewares ...Middleware) {
	previousGroup := hr.currentGroup
	newGroup := &RouteGroup{
		prefix:      prefix,
		middlewares: middlewares,
	}

	if previousGroup != nil {
		newGroup.prefix = path.Join(previousGroup.prefix, newGroup.prefix)
		newGroup.middlewares = append(previousGroup.middlewares, newGroup.middlewares...)
	}

	hr.currentGroup = newGroup
	fn(hr)
	hr.currentGroup = previousGroup
}

func (hr *Router) appendPatternAndMiddlewares(pattern string, middlewares []Middleware) (string, []Middleware) {
	fullPattern := pattern
	fullMiddlewares := middlewares

	if hr.currentGroup != nil {
		fullPattern = path.Join(hr.currentGroup.prefix, pattern)
		fullMiddlewares = append(hr.currentGroup.middlewares, middlewares...)
	}
	return fullPattern, fullMiddlewares
}

func (hr *Router) registerMiddlewares(app *App) {
	for _, plugin := range app.plugins {
		hr.httpMiddlewares = append(hr.httpMiddlewares, plugin.Middlewares()...)
	}
	for _, mw := range hr.httpMiddlewares {
		hr.Use(adaptMiddleware(app, mw))
	}
	//hr.Use(hr.httpMiddlewares...)
}

func (hr *Router) registerRoutes(app *App) {
	hr.routeRegistrar(hr)
	for _, plugin := range app.plugins {
		for _, route := range plugin.Routes() {
			hr.routes = append(hr.routes, route)
		}
	}

	for _, route := range hr.routes {
		hr.MethodFunc(route.Method, route.Path, makeHandlerFunc(app, route.Handler, route.Middlewares...))
	}
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

func makeHandlerFunc(app *App, handler Handler, middlewares ...Middleware) http.HandlerFunc {
	finalHandler := handler
	for _, middleware := range middlewares {
		finalHandler = middleware(finalHandler)
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{sync.Mutex{}, app, container.New(), r, w}
		if !app.isContextReady {
			app.isContextReady = true
		}
		if err := finalHandler(ctx); err != nil {
			logger.V().Error(err.Error())
			if !ctx.WantsJSON() {
				ctx.Redirect(http.StatusFound, "/error")
				return
			}
			ctx.JSON(http.StatusInternalServerError, M{"error": err.Error()})
		}
		return
	}

	return fn

	//return app.I.Middleware(http.HandlerFunc(fn)).ServeHTTP
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
