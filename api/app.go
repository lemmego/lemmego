package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"reflect"
	"sync"
	"syscall"

	"pressebo/api/db"
	"pressebo/api/logger"

	"github.com/go-chi/chi/v5"
	"github.com/golobby/container/v3"
	"github.com/joho/godotenv"
	inertia "github.com/romsar/gonertia"
	"github.com/spf13/cobra"
)

type PluginID string

type PluginRegistry map[PluginID]Plugin

// Find a plugin
func (r PluginRegistry) Find(namespace string) Plugin {
	return r[PluginID(namespace)]
}

// Add a plugin
func (r PluginRegistry) Add(plugin Plugin) {
	r[PluginID(plugin.Namespace())] = plugin
}

type M map[string]any

type Handler func(c *Context) error
type Middleware func(next Handler) Handler

type AppConfig struct {
	DbConfig    *db.DBConfig
	AppName     string
	AppPort     int
	TemplateDir string
}

type Plugin interface {
	Boot(a *App) error
	Commands() []*cobra.Command
	Namespace() string
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
	session        *Session
	// config           *AppConfig
	config           ConfigMap
	plugins          PluginRegistry
	services         []ServiceProvider
	routeMiddlewares map[string]Middleware
	hooks            *AppHooks
	router           Router
	db               db.DBSession
	dbFunc           func(config *db.DBConfig) (db.DBSession, error)
	I                *inertia.Inertia
}

type Options struct {
	container.Container
	*Session
	Config           ConfigMap
	Plugins          map[PluginID]Plugin
	Providers        []ServiceProvider
	routeMiddlewares map[string]Middleware
	Hooks            *AppHooks
	Router           Router
}

type OptFunc func(opts *Options)

func (a *App) Plugin(namespace string) Plugin {
	return a.plugins.Find(namespace)
}

func (a *App) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, m := range middlewares {
		a.router.Use(m)
	}
}

func (a *App) Get(pattern string, handler Handler, middlewares ...Middleware) {
	a.router.MethodFunc(http.MethodGet, pattern, makeHandlerFunc(a, handler, middlewares...))
}

func (a *App) Post(pattern string, handler Handler, middlewares ...Middleware) {
	a.router.MethodFunc(http.MethodPost, pattern, makeHandlerFunc(a, handler, middlewares...))
}

func (a *App) Put(pattern string, handler Handler, middlewares ...Middleware) {
	a.router.MethodFunc(http.MethodPut, pattern, makeHandlerFunc(a, handler, middlewares...))
}

func (a *App) Patch(pattern string, handler Handler, middlewares ...Middleware) {
	a.router.MethodFunc(http.MethodPatch, pattern, makeHandlerFunc(a, handler, middlewares...))
}

func (a *App) Delete(pattern string, handler Handler, middlewares ...Middleware) {
	a.router.MethodFunc(http.MethodDelete, pattern, makeHandlerFunc(a, handler, middlewares...))
}

func (a *App) Router() chi.Router {
	return a.router
}

func (a *App) Session() *Session {
	return a.session
}

func (a *App) Db() db.DBSession {
	return a.db
}

// func (a *App) Inertia() *inertia.Inertia {
// 	return a.inertia
// }

func (a *App) DbFunc(config *db.DBConfig) (db.DBSession, error) {
	return a.dbFunc(config)
}

func getDefaultConfig() ConfigMap {
	return Conf
	// host := MustEnv("DB_HOST", "localhost")
	// port := MustEnv("DB_PORT", 5432)
	// database := MustEnv("DB_DATABASE", "fluentapp")
	// user := MustEnv("DB_USERNAME", "fluentapp")
	// password := MustEnv("DB_PASSWORD", "fluentapp")
	// appName := MustEnv("APP_NAME", "FluentApp")
	// appPort := MustEnv("APP_PORT", 3000)

	// return &AppConfig{
	// 	AppName: appName,
	// 	AppPort: appPort,
	// 	DbConfig: &DBConfig{
	// 		Host:     host,
	// 		Port:     port,
	// 		Database: database,
	// 		User:     user,
	// 		Password: password,
	// 	},
	// 	TemplateDir: "templates",
	// }
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
		NewRouter(),
	}
}

func WithPlugins(plugins map[PluginID]Plugin) OptFunc {
	return func(opts *Options) {
		opts.Plugins = plugins
	}
}

func WithProviders(providers []ServiceProvider) OptFunc {
	return func(opts *Options) {
		opts.Providers = providers
	}
}

func WithHooks(hooks *AppHooks) OptFunc {
	return func(opts *Options) {
		opts.Hooks = hooks
	}
}

func WithRouter(router Router) OptFunc {
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
		// Copy template files listed in the Views() method to the app's template directory
		for name, content := range plugin.Templates() {
			filePath := filepath.Join(opts.Config.get("app.templateDir").(string), name)
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

		// if err := plugin.Boot(opts.Container); err != nil {
		// 	panic(err)
		// }

	}

	i := initInertia()

	app := &App{
		// opts.Container,
		false,
		sync.Mutex{},
		nil,
		opts.Config,
		opts.Plugins,
		opts.Providers,
		routeMiddlewares,
		opts.Hooks,
		opts.Router,
		nil,
		nil,
		i,
	}
	return app
}

// func (app *App) Container() container.Container {
// 	return app.container
// }

func (app *App) registerServices(services []ServiceProvider) {
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

func (app *App) registerMiddlewares(middlewares []func(http.Handler) http.Handler) {
	for _, middleware := range middlewares {
		app.router.Use(middleware)
	}
	for _, plugin := range app.plugins {
		for _, middleware := range plugin.Middlewares() {
			app.router.Use(middleware)
		}
	}
}

func (app *App) registerRoutes() {
	// for _, route := range routes {
	// 	app.Router.MethodFunc(route.Method, route.Path, makeHandlerFunc(app, route.Handler, route.Middlewares...))
	// }

	for _, plugin := range app.plugins {
		for _, route := range plugin.Routes() {
			app.router.MethodFunc(route.Method, route.Path, makeHandlerFunc(app, route.Handler, route.Middlewares...))
		}
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
			logger.Log().Error(err.Error())
			if !ctx.WantsJSON() {
				ctx.Redirect(http.StatusFound, "/error")
				return
			}
			ctx.JSON(http.StatusInternalServerError, M{"error": err.Error()})
		}
		return
	}

	return http.HandlerFunc(fn)
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

func (a *App) Run() {
	a.registerServices([]ServiceProvider{
		&DatabaseServiceProvider{},
		&SessionServiceProvider{},
		&AuthServiceProvider{},
	})

	a.registerRoutes()

	a.router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	a.router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	for _, plugin := range a.plugins {
		if err := plugin.Boot(a); err != nil {
			panic(err)
		}
	}

	fmt.Println(fmt.Sprintf("%s is running on port %d...", a.config.get("app.name"), a.config.get("app.port")))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.config.get("app.port")), a.session.LoadAndSave(a.router)); err != nil {
		panic(err)
	}
}

func (a *App) HandleSignals() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	sig := <-signalChannel
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		a.Shutdown()
		os.Exit(0)
	}
}

func (a *App) Shutdown() {
	log.Println("Shutting down application...")
	sessName := a.db.Name()
	err := a.db.Close()
	if err != nil {
		log.Fatal("Error closing database connection:", err)
	}
	log.Println("Database connection", sessName, "closed.")
}
