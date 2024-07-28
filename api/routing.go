package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/core"
	inertia "github.com/romsar/gonertia"
	"lemmego/api/logger"
	"lemmego/api/vee"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"sync"
)

const HTTPInKey = "input"

type Handler func(c *Context) error

type Middleware func(next Handler) Handler

type HTTPMiddleware func(http.Handler) http.Handler

type RouteRegistrarFunc func(r *Router)

type InertiaFlashProvider struct {
	errors map[string]inertia.ValidationErrors
}

func NewInertiaFlashProvider() *InertiaFlashProvider {
	return &InertiaFlashProvider{errors: make(map[string]inertia.ValidationErrors)}
}

func (p *InertiaFlashProvider) FlashErrors(ctx context.Context, errors inertia.ValidationErrors) error {
	if sessionID, ok := ctx.Value("sessionID").(string); ok {
		p.errors[sessionID] = errors
	}
	return nil
}

func (p *InertiaFlashProvider) GetErrors(ctx context.Context) (inertia.ValidationErrors, error) {
	var inertiaErrors inertia.ValidationErrors
	if sessionID, ok := ctx.Value("sessionID").(string); ok {
		inertiaErrors = p.errors[sessionID]
		p.errors[sessionID] = nil
	}
	return inertiaErrors, nil
}

type Router struct {
	routes           []*Route
	routeMiddlewares map[string]Middleware
	httpMiddlewares  []HTTPMiddleware
	currentGroup     *RouteGroup
	basePrefix       string
	mux              *http.ServeMux
}

// NewRouter creates a new HTTPRouter-based router
func NewRouter() *Router {
	return &Router{
		routes:           []*Route{},
		routeMiddlewares: make(map[string]Middleware),
		httpMiddlewares:  []HTTPMiddleware{},
		mux:              http.NewServeMux(),
	}
}

func (r *Router) HasRoute(method string, pattern string) bool {
	return slices.ContainsFunc(r.routes, func(route *Route) bool {
		return route.Method == method && route.Path == pattern
	})
}

func (r *Router) setRouteMiddleware(middlewares map[string]Middleware) {
	r.routeMiddlewares = middlewares
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.Handler = r.mux
	for i := len(r.httpMiddlewares) - 1; i >= 0; i-- {
		handler = r.httpMiddlewares[i](handler)
	}
	handler.ServeHTTP(w, req)
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(pattern, handler)
}

func (r *Router) Get(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodGet, pattern, handler)
}

func (r *Router) Post(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodPost, pattern, handler)
}

func (r *Router) Put(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodPut, pattern, handler)
}

func (r *Router) Patch(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodPatch, pattern, handler)
}

func (r *Router) Delete(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodDelete, pattern, handler)
}

func (r *Router) Connect(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodConnect, pattern, handler)
}

func (r *Router) Head(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodHead, pattern, handler)
}

func (r *Router) Options(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodOptions, pattern, handler)
}

func (r *Router) Trace(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodTrace, pattern, handler)
}

func (r *Router) addRoute(method, pattern string, handler Handler) *Route {
	fullPath := pattern
	if r.currentGroup != nil {
		fullPath = path.Join(r.currentGroup.prefix, pattern)
	}
	route := &Route{Method: method, Path: fullPath, Handler: handler}
	if r.currentGroup != nil {
		r.currentGroup.routes = append(r.currentGroup.routes, route)
	}
	r.routes = append(r.routes, route)
	return route
}

// RouteGroup represents a group of routes
type RouteGroup struct {
	prefix           string
	beforeMiddleware []Middleware
	afterMiddleware  []Middleware
	routes           []*Route
}

// Use adds one or more standard net/http middleware to the router
func (r *Router) Use(middlewares ...HTTPMiddleware) {
	r.httpMiddlewares = append(r.httpMiddlewares, middlewares...)
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
func (r *Router) Group(prefix string, fn func(r *Router)) *RouteGroup {
	previousGroup := r.currentGroup
	newGroup := &RouteGroup{
		prefix: prefix,
	}

	if previousGroup != nil {
		newGroup.prefix = path.Join(previousGroup.prefix, newGroup.prefix)
	}

	r.currentGroup = newGroup

	fn(r)

	r.currentGroup = previousGroup

	return newGroup
}

func (r *Router) registerMiddlewares(app *App) {
	for _, plugin := range app.plugins {
		for _, mw := range plugin.Middlewares() {
			r.Use(mw)
		}
	}
}

func (r *Router) registerRoutes(app *App) {
	for _, plugin := range app.plugins {
		for _, route := range plugin.Routes() {
			if !r.HasRoute(route.Method, route.Path) {
				r.routes = append(r.routes, route)
			}
		}
	}

	for _, route := range r.routes {
		handler := route.Handler

		// Apply group middlewares
		if r.currentGroup != nil {
			for i := len(r.currentGroup.beforeMiddleware) - 1; i >= 0; i-- {
				handler = r.currentGroup.beforeMiddleware[i](handler)
			}
			for i := len(r.currentGroup.afterMiddleware) - 1; i >= 0; i-- {
				nextHandler := handler
				handler = func(c *Context) error {
					err := nextHandler(c)
					if err != nil {
						return err
					}
					return r.currentGroup.afterMiddleware[i](func(*Context) error { return nil })(c)
				}
			}
		}

		r.mux.HandleFunc(route.Method+" "+route.Path, makeHandlerFunc(app, &Route{
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

func AdaptMiddleware(app *App, m Middleware) Middleware {
	return func(next Handler) Handler {
		return func(c *Context) error {
			wrappedNext := func(c *Context) error {
				return next(c)
			}

			wrappedHandler := m(wrappedNext)

			return wrappedHandler(c)
		}
	}
}

func makeHandlerFunc(app *App, route *Route) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//app.Reset()
		token := app.Session().Token(r.Context())
		if token != "" {
			r = r.WithContext(context.WithValue(r.Context(), "sessionID", token))
			log.Println("Current SessionID: ", token)
		}
		ctx := &Context{
			Mutex:          sync.Mutex{},
			app:            app,
			request:        r,
			responseWriter: w,
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
			if errors.As(err, &vee.Errors{}) {
				ctx.ValidationError(err)
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
		inertia.WithVersionFromFile(manifestPath),
		inertia.WithSSR(),
		//inertia.WithVersion("1.0"),
		inertia.WithFlashProvider(NewInertiaFlashProvider()),
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
//				ctx.Set(HTTPInKey, input)
//
//				// Create a new request with the context
//				r = r.WithContext(context.WithValue(r.Context(), HTTPInKey, ctx))
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

			ctx.Set(HTTPInKey, input)
			return next(ctx)
		}
	}
}
