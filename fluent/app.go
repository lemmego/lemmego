//go:generate go -version
package fluent

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/golobby/container/v3"
	"github.com/joho/godotenv"
)

type M map[string]any

type Handler func(c *Context) error
type Middleware func(next Handler) Handler

type Config struct {
	DbConfig    *DBConfig
	AppName     string
	AppPort     int
	TemplateDir string
}

// type App interface {
// 	GetConfig() *Config
// }

type Plugin interface {
	Namespace() string
	Init(c container.Container) error
	EventListeners() map[string]func()
	Migrations() []string
	Templates() map[string][]byte
	Middlewares() []func(http.Handler) http.Handler
	RouteMiddlewares() map[string]Middleware
	Routes() []*Route
	Webhooks() []string
}

type AppHooks struct {
	BeforeStart func()
	AfterStart  func()
}

type App struct {
	isContextReady bool
	mu             sync.Mutex
	container.Container
	session          *Session
	Config           *Config
	plugins          []Plugin
	services         []ServiceProvider
	routeMiddlewares map[string]Middleware
	Hooks            *AppHooks
	Router           chi.Router
}

type Options struct {
	container.Container
	*Session
	Config           *Config
	Plugins          []Plugin
	Services         []ServiceProvider
	routeMiddlewares map[string]Middleware
	Hooks            *AppHooks
	Router           chi.Router
}

type OptFunc func(opts *Options)

func (a *App) Get(pattern string, handler Handler, middlewares ...Middleware) (*Route, error) {
	if a.isContextReady {
		return nil, fmt.Errorf("cannot add route after context is ready")
	}
	a.Router.MethodFunc(http.MethodGet, pattern, makeHandlerFunc(a, handler, middlewares...))
	return NewRoute(http.MethodGet, pattern, handler, middlewares...), nil
}

func (a *App) Post(pattern string, handler Handler, middlewares ...Middleware) (*Route, error) {
	if a.isContextReady {
		return nil, fmt.Errorf("cannot add route after context is ready")
	}
	a.Router.MethodFunc(http.MethodPost, pattern, makeHandlerFunc(a, handler, middlewares...))
	return NewRoute(http.MethodPost, pattern, handler, middlewares...), nil
}

func (a *App) Put(pattern string, handler Handler, middlewares ...Middleware) (*Route, error) {
	if a.isContextReady {
		return nil, fmt.Errorf("cannot add route after context is ready")
	}
	a.Router.MethodFunc(http.MethodPut, pattern, makeHandlerFunc(a, handler, middlewares...))
	return NewRoute(http.MethodPut, pattern, handler, middlewares...), nil
}

func (a *App) Patch(pattern string, handler Handler, middlewares ...Middleware) (*Route, error) {
	if a.isContextReady {
		return nil, fmt.Errorf("cannot add route after context is ready")
	}
	a.Router.MethodFunc(http.MethodPatch, pattern, makeHandlerFunc(a, handler, middlewares...))
	return NewRoute(http.MethodPatch, pattern, handler, middlewares...), nil
}

func (a *App) Delete(pattern string, handler Handler, middlewares ...Middleware) (*Route, error) {
	if a.isContextReady {
		return nil, fmt.Errorf("cannot add route after context is ready")
	}
	a.Router.MethodFunc(http.MethodDelete, pattern, makeHandlerFunc(a, handler, middlewares...))
	return NewRoute(http.MethodDelete, pattern, handler, middlewares...), nil
}

func getDefaultConfig() *Config {
	var appName, host, database, user, password string
	var appPort, port int

	if val, ok := os.LookupEnv("DB_HOST"); ok {
		host = val
	} else {
		host = "localhost"
	}

	if val, ok := os.LookupEnv("DB_PORT"); ok {
		port, _ = strconv.Atoi(val)
	} else {
		port = 5432
	}

	if val, ok := os.LookupEnv("DB_DATABASE"); ok {
		database = val
	} else {
		database = "fluentapp"
	}

	if val, ok := os.LookupEnv("DB_USERNAME"); ok {
		user = val
	} else {
		user = "fluentapp"
	}

	if val, ok := os.LookupEnv("DB_PASSWORD"); ok {
		password = val
	} else {
		password = "fluentapp"
	}

	if val, ok := os.LookupEnv("APP_NAME"); ok {
		appName = val
	} else {
		appName = "FluentApp"
	}

	if val, ok := os.LookupEnv("APP_PORT"); ok {
		appPort, _ = strconv.Atoi(val)
	} else {
		appPort = 3000
	}

	return &Config{
		AppName: appName,
		AppPort: appPort,
		DbConfig: &DBConfig{
			Host:     host,
			Port:     port,
			Database: database,
			User:     user,
			Password: password,
		},
		TemplateDir: "templates",
	}
}

func defaultOptions() *Options {
	return &Options{
		container.New(),
		nil,
		getDefaultConfig(),
		nil,
		nil,
		nil,
		nil,
		chi.NewRouter(),
	}
}

func WithConfig(config *Config) OptFunc {
	return func(opts *Options) {
		opts.Config = config
	}
}

func WithPlugins(plugins []Plugin) OptFunc {
	return func(opts *Options) {
		opts.Plugins = plugins
	}
}

func WithServices(services []ServiceProvider) OptFunc {
	return func(opts *Options) {
		opts.Services = services
	}
}

func WithHooks(hooks *AppHooks) OptFunc {
	return func(opts *Options) {
		opts.Hooks = hooks
	}
}

func WithRouter(router chi.Router) OptFunc {
	return func(opts *Options) {
		opts.Router = router
	}
}

func WithContainer(container container.Container) OptFunc {
	return func(opts *Options) {
		opts.Container = container
	}
}

func WithSession(sm *Session) OptFunc {
	return func(opts *Options) {
		opts.Session = sm
	}
}

func NewApp(options ...OptFunc) *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	opts := defaultOptions()
	for _, option := range options {
		option(opts)
	}

	// Check if plugins have duplicate namespaces
	var namespaces []string

	for _, plugin := range opts.Plugins {
		if plugin.Namespace() == "" {
			panic("Plugin namespace cannot be empty. Please set a namespace for plugin: " + reflect.TypeOf(plugin).String())
		}
		for _, namespace := range namespaces {
			if namespace == plugin.Namespace() {
				panic("Duplicate plugin namespace: " + namespace)
			}
		}
		namespaces = append(namespaces, plugin.Namespace())
	}

	var routeMiddlewares map[string]Middleware

	for _, plugin := range opts.Plugins {
		if err := plugin.Init(opts.Container); err != nil {
			panic(err)
		}

		// Copy template files listed in the Views() method to the app's template directory
		for name, content := range plugin.Templates() {
			filePath := filepath.Join(opts.Config.TemplateDir, name)
			if _, err := os.Stat(filePath); err != nil {
				ioutil.WriteFile(filePath, []byte(content), 0644)
			}
		}

		for middlewareName, middleware := range plugin.RouteMiddlewares() {
			if routeMiddlewares == nil {
				routeMiddlewares = make(map[string]Middleware)
			}
			key := plugin.Namespace() + "." + middlewareName
			if _, ok := routeMiddlewares[key]; ok {
				panic(fmt.Sprintf("Middleware %s already registered", plugin.Namespace()+"."+middlewareName))
			}
			routeMiddlewares[key] = middleware
		}
	}

	return &App{
		false,
		sync.Mutex{},
		opts.Container,
		nil,
		opts.Config,
		opts.Plugins,
		opts.Services,
		routeMiddlewares,
		opts.Hooks,
		opts.Router,
	}
}

func (app *App) RegisterServices(services []ServiceProvider) {
	for _, svc := range services {
		if reflect.TypeOf(svc).Kind() != reflect.Ptr {
			panic("Service must be a pointer")
		}
		if reflect.TypeOf(svc).Elem().Kind() != reflect.Struct {
			panic("Service must be a struct")
		}
		if reflect.TypeOf(svc).Elem().Field(0).Name != "BaseServiceProvider" {
			panic("Service must extend BaseServiceProvider")
		}
		if reflect.TypeOf(svc).Elem().Field(0).Type != reflect.TypeOf(BaseServiceProvider{}) {
			panic("Service must extend BaseServiceProvider")
		}

		// Check if service implements ServiceProvider interface, not necessary if type hinted
		if reflect.TypeOf(svc).Implements(reflect.TypeOf((*ServiceProvider)(nil)).Elem()) == true {
			println("Registering service: " + reflect.TypeOf(svc).Elem().Name())
			svc.(ServiceProvider).Register(app)
			app.services = append(app.services, svc.(ServiceProvider))
		} else {
			panic("Service must implement ServiceProvider interface")
		}

	}
	for _, service := range services {
		if reflect.TypeOf(service).Implements(reflect.TypeOf((*ServiceProvider)(nil)).Elem()) {
			service.(ServiceProvider).Boot()
		}
	}
}

func (app *App) RegisterMiddlewares(middlewares []func(http.Handler) http.Handler) {
	for _, middleware := range middlewares {
		app.Router.Use(middleware)
	}
	for _, plugin := range app.plugins {
		for _, middleware := range plugin.Middlewares() {
			app.Router.Use(middleware)
		}
	}
}

func (app *App) RegisterRoutes() {
	// for _, route := range routes {
	// 	app.Router.MethodFunc(route.Method, route.Path, makeHandlerFunc(app, route.Handler, route.Middlewares...))
	// }

	for _, plugin := range app.plugins {
		for _, route := range plugin.Routes() {
			app.Router.MethodFunc(route.Method, route.Path, makeHandlerFunc(app, route.Handler))
		}
	}
}

func makeHandlerFunc(app *App, handler Handler, middlewares ...Middleware) http.HandlerFunc {
	finalHandler := handler
	for _, middleware := range middlewares {
		finalHandler = middleware(finalHandler)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{app, r, w}
		if !app.isContextReady {
			app.isContextReady = true
		}
		if err := finalHandler(ctx); err != nil {
			log.Println(err)
			if !ctx.WantsJSON() {
				ctx.Redirect("/error", http.StatusFound)
				return
			}
			ctx.JSON(http.StatusInternalServerError, M{"error": err.Error()})
		}
		return
	}
}

func (a *App) Run() {
	a.RegisterServices([]ServiceProvider{
		&DatabaseServiceProvider{},
		&SessionServiceProvider{},
		&AuthServiceProvider{},
	})

	a.RegisterRoutes()

	a.Router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println(fmt.Sprintf("%s is running on port %d...", a.Config.AppName, a.Config.AppPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.Config.AppPort), a.session.LoadAndSave(a.Router)); err != nil {
		panic(err)
	}
}
