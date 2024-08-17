package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"slices"

	"github.com/lemmego/lemmego/api/config"

	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/core"
	inertia "github.com/romsar/gonertia"
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
	router           *Router
	prefix           string
	beforeMiddleware []Handler
	afterMiddleware  []Handler
}

func (g *Group) Group(prefix string) *Group {
	return &Group{
		router:           g.router,
		prefix:           path.Join(g.prefix, prefix),
		beforeMiddleware: append([]Handler{}, g.beforeMiddleware...),
		afterMiddleware:  append([]Handler{}, g.afterMiddleware...),
	}
}

func (g *Group) UseBefore(handlers ...Handler) {
	g.beforeMiddleware = append(g.beforeMiddleware, handlers...)
}

func (g *Group) UseAfter(handlers ...Handler) {
	g.afterMiddleware = append(handlers, g.afterMiddleware...)
}

func (g *Group) addRoute(method, pattern string, handlers ...Handler) *Route {
	fullPath := path.Join(g.prefix, pattern)
	route := &Route{
		Method:           method,
		Path:             fullPath,
		Handlers:         handlers,
		BeforeMiddleware: append(append([]Handler{}, g.router.beforeMiddleware...), g.beforeMiddleware...),
		AfterMiddleware:  append(append([]Handler{}, g.afterMiddleware...), g.router.afterMiddleware...),
		router:           g.router,
	}
	g.router.routes = append(g.router.routes, route)
	return route
}

func (g *Group) Get(pattern string, handlers ...Handler) *Route {
	return g.addRoute(http.MethodGet, pattern, handlers...)
}

func (g *Group) Post(pattern string, handlers ...Handler) *Route {
	return g.addRoute(http.MethodPost, pattern, handlers...)
}

func (g *Group) Put(pattern string, handlers ...Handler) *Route {
	return g.addRoute(http.MethodPut, pattern, handlers...)
}

func (g *Group) Patch(pattern string, handlers ...Handler) *Route {
	return g.addRoute(http.MethodPatch, pattern, handlers...)
}

func (g *Group) Delete(pattern string, handlers ...Handler) *Route {
	return g.addRoute(http.MethodDelete, pattern, handlers...)
}

type Router struct {
	routes           []*Route
	httpMiddlewares  []HTTPMiddleware
	basePrefix       string
	mux              *http.ServeMux
	beforeMiddleware []Handler
	afterMiddleware  []Handler

	namedMiddleware map[string]Handler
}

// NewRouter creates a new HTTPRouter-based router
func NewRouter() *Router {
	return &Router{
		routes:           []*Route{},
		httpMiddlewares:  []HTTPMiddleware{},
		mux:              http.NewServeMux(),
		beforeMiddleware: []Handler{},
		afterMiddleware:  []Handler{},

		namedMiddleware: make(map[string]Handler),
	}
}

func (r *Router) Group(prefix string) *Group {
	return &Group{
		router:           r,
		prefix:           prefix,
		beforeMiddleware: []Handler{},
		afterMiddleware:  []Handler{},
	}
}

func (r *Router) UseBefore(handlers ...Handler) {
	r.beforeMiddleware = append(r.beforeMiddleware, handlers...)
}

func (r *Router) UseAfter(handlers ...Handler) {
	r.afterMiddleware = append(handlers, r.afterMiddleware...)
}

func (r *Router) HasRoute(method string, pattern string) bool {
	return slices.ContainsFunc(r.routes, func(route *Route) bool {
		return route.Method == method && route.Path == pattern
	})
}

func (r *Router) setNamedMiddleware(middlewares map[string]Handler) {
	r.namedMiddleware = middlewares
}

func (r *Router) Mw(name string) Handler {
	if val, ok := r.namedMiddleware[name]; ok {
		return val
	}

	return func(c *Context) error {
		return c.Next()
	}
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

func (r *Router) Get(pattern string, handlers ...Handler) *Route {
	return r.addRoute(http.MethodGet, pattern, handlers...)
}

func (r *Router) Post(pattern string, handlers ...Handler) *Route {
	return r.addRoute(http.MethodPost, pattern, handlers...)
}

func (r *Router) Put(pattern string, handlers ...Handler) *Route {
	return r.addRoute(http.MethodPut, pattern, handlers...)
}

func (r *Router) Patch(pattern string, handlers ...Handler) *Route {
	return r.addRoute(http.MethodPatch, pattern, handlers...)
}

func (r *Router) Delete(pattern string, handlers ...Handler) *Route {
	return r.addRoute(http.MethodDelete, pattern, handlers...)
}

func (r *Router) Connect(pattern string, handlers ...Handler) *Route {
	return r.addRoute(http.MethodConnect, pattern, handlers...)
}

func (r *Router) Head(pattern string, handlers ...Handler) *Route {
	return r.addRoute(http.MethodHead, pattern, handlers...)
}

func (r *Router) Options(pattern string, handlers ...Handler) *Route {
	return r.addRoute(http.MethodOptions, pattern, handlers...)
}

func (r *Router) Trace(pattern string, handlers ...Handler) *Route {
	return r.addRoute(http.MethodTrace, pattern, handlers...)
}

func (r *Router) addRoute(method, pattern string, handlers ...Handler) *Route {
	fullPath := path.Join(r.basePrefix, pattern)
	route := &Route{
		Method:           method,
		Path:             fullPath,
		Handlers:         handlers,
		BeforeMiddleware: r.beforeMiddleware,
		AfterMiddleware:  r.afterMiddleware,
		router:           r,
	}
	r.routes = append(r.routes, route)
	log.Printf("Added route: %s %s, router: %p", method, fullPath, r)
	return route
}

// Use adds one or more standard net/http middleware to the router
func (r *Router) Use(middlewares ...HTTPMiddleware) {
	r.httpMiddlewares = append(r.httpMiddlewares, middlewares...)
}

type Route struct {
	Method           string
	Path             string
	Handlers         []Handler
	BeforeMiddleware []Handler
	AfterMiddleware  []Handler
	router           *Router
}

func (r *Route) UseBefore(handlers ...Handler) *Route {
	r.BeforeMiddleware = append(r.BeforeMiddleware, handlers...)
	return r
}

func (r *Route) UseAfter(handlers ...Handler) *Route {
	r.AfterMiddleware = append(handlers, r.AfterMiddleware...)
	return r
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
	i.ShareTemplateData("env", config.Config("app.env").(string))

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
