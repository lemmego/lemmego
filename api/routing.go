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

type Group struct {
	router            *Router
	prefix            string
	beforeMiddlewares []Middleware
	afterMiddlewares  []Middleware
}

func (g *Group) Group(prefix string) *Group {
	return &Group{
		router:            g.router,
		prefix:            path.Join(g.prefix, prefix),
		beforeMiddlewares: append([]Middleware{}, g.beforeMiddlewares...),
		afterMiddlewares:  append([]Middleware{}, g.afterMiddlewares...),
	}
}

func (g *Group) UseBefore(middleware ...Middleware) {
	g.beforeMiddlewares = append(g.beforeMiddlewares, middleware...)
}

func (g *Group) UseAfter(middleware ...Middleware) {
	g.afterMiddlewares = append(g.afterMiddlewares, middleware...)
}

func (g *Group) addRoute(method, pattern string, handler Handler) *Route {
	fullPath := path.Join(g.prefix, pattern)
	route := &Route{
		Method:           method,
		Path:             fullPath,
		Handler:          handler,
		router:           g.router,
		BeforeMiddleware: append([]Middleware{}, g.beforeMiddlewares...),
		AfterMiddleware:  append([]Middleware{}, g.afterMiddlewares...),
	}
	g.router.routes = append(g.router.routes, route)
	log.Printf("Added route to group: %s %s, router: %p", method, fullPath, g.router)
	return route
}

func (g *Group) Get(pattern string, handler Handler) *Route {
	return g.addRoute(http.MethodGet, pattern, handler)
}

func (g *Group) Post(pattern string, handler Handler) *Route {
	return g.addRoute(http.MethodPost, pattern, handler)
}

func (g *Group) Put(pattern string, handler Handler) *Route {
	return g.addRoute(http.MethodPut, pattern, handler)
}

func (g *Group) Patch(pattern string, handler Handler) *Route {
	return g.addRoute(http.MethodPatch, pattern, handler)
}

func (g *Group) Delete(pattern string, handler Handler) *Route {
	return g.addRoute(http.MethodDelete, pattern, handler)
}

type Router struct {
	routes            []*Route
	routeMiddlewares  map[string]Middleware
	httpMiddlewares   []HTTPMiddleware
	basePrefix        string
	mux               *http.ServeMux
	beforeMiddlewares []Middleware
	afterMiddlewares  []Middleware
}

// NewRouter creates a new HTTPRouter-based router
func NewRouter() *Router {
	return &Router{
		routes:            []*Route{},
		routeMiddlewares:  make(map[string]Middleware),
		httpMiddlewares:   []HTTPMiddleware{},
		mux:               http.NewServeMux(),
		beforeMiddlewares: []Middleware{},
		afterMiddlewares:  []Middleware{},
	}
}

func (r *Router) Group(prefix string) *Group {
	group := &Group{
		router:            r,
		prefix:            prefix,
		beforeMiddlewares: []Middleware{},
		afterMiddlewares:  []Middleware{},
	}
	log.Printf("Created group with prefix: %s, router: %p", prefix, r)
	return group
}

func (r *Router) UseBefore(middleware ...Middleware) {
	r.beforeMiddlewares = append(r.beforeMiddlewares, middleware...)
}

func (r *Router) UseAfter(middleware ...Middleware) {
	r.afterMiddlewares = append(r.afterMiddlewares, middleware...)
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
	fullPath := path.Join(r.basePrefix, pattern)
	route := &Route{
		Method:           method,
		Path:             fullPath,
		Handler:          handler,
		router:           r,
		BeforeMiddleware: []Middleware{},
		AfterMiddleware:  []Middleware{},
	}
	r.routes = append(r.routes, route)
	log.Printf("Added route: %s %s, router: %p", method, fullPath, r)
	return route
}

// Use adds one or more standard net/http middleware to the router
func (r *Router) Use(middlewares ...HTTPMiddleware) {
	r.httpMiddlewares = append(r.httpMiddlewares, middlewares...)
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
				log.Println("Adding route for the", plugin.Namespace(), "plugin:", route.Method, route.Path)
				r.addRoute(route.Method, route.Path, route.Handler)
				//r.routes = append(r.routes, route)
			}
		}
	}

	for _, route := range r.routes {
		log.Printf("Registering route: %s %s, router: %p", route.Method, route.Path, route.router)
		r.mux.HandleFunc(route.Method+" "+route.Path, func(w http.ResponseWriter, req *http.Request) {
			makeHandlerFunc(app, route, r)(w, req)
		})
	}
}

type Route struct {
	Method           string
	Path             string
	Handler          Handler
	AfterMiddleware  []Middleware
	BeforeMiddleware []Middleware
	router           *Router
}

func (r *Route) UseBefore(middleware ...Middleware) *Route {
	r.BeforeMiddleware = append(r.BeforeMiddleware, middleware...)
	return r
}

func (r *Route) UseAfter(middleware ...Middleware) *Route {
	r.AfterMiddleware = append(r.AfterMiddleware, middleware...)
	return r
}

func makeHandlerFunc(app *App, route *Route, router *Router) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling request for route: %s %s, router: %p", route.Method, route.Path, router)
		if route.router == nil {
			log.Printf("WARNING: route.router is nil for %s %s", route.Method, route.Path)
			// Handle the error condition, maybe return a 500 Internal Server Error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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
		log.Printf("Applying middlewares for route: %s %s", route.Method, route.Path)

		// Apply router-level after middlewares
		for i := len(router.afterMiddlewares) - 1; i >= 0; i-- {
			afterMiddleware := router.afterMiddlewares[i]
			nextChain := chain
			chain = func(c *Context) error {
				err := nextChain(c)
				if err != nil {
					return err
				}
				return afterMiddleware(func(*Context) error { return nil })(c)
			}
		}

		// Apply route-specific after middlewares
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

		// Apply route-specific before middlewares
		for i := len(route.BeforeMiddleware) - 1; i >= 0; i-- {
			chain = route.BeforeMiddleware[i](chain)
		}

		// Apply router-level before middlewares
		for i := len(route.router.beforeMiddlewares) - 1; i >= 0; i-- {
			chain = route.router.beforeMiddlewares[i](chain)
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
