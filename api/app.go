package api

import (
	"context"
	"fmt"
	"lemmego/api/fsys"
	"lemmego/api/session"
	// "lemmego/api/session"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"sync"
	"syscall"

	"github.com/golobby/container/v3"
	inertia "github.com/romsar/gonertia"
	"github.com/spf13/cobra"
	"lemmego/api/db"
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

//type Middleware func(c *Context) func(next http.Handler) http.Handler

//type Handler[T any] func(c T) error

//type Middleware[T any] func(next Handler[T]) Handler[T]

type AppConfig struct {
	DbConfig    *db.DBConfig
	AppName     string
	AppPort     int
	TemplateDir string
}

type Plugin interface {
	Boot(a *App) error
	InstallCommand() *cobra.Command
	Commands() []*cobra.Command
	Namespace() string
	EventListeners() map[string]func()
	Migrations() []string
	Templates() map[string][]byte
	Middlewares() []Middleware
	RouteMiddlewares() map[string]Middleware
	Routes() []*Route
	Webhooks() []string
	// Generators() []cli.Generator
}

type AppHooks struct {
	BeforeStart func()
	AfterStart  func()
}

type App struct {
	isContextReady bool
	mu             sync.Mutex
	session        *session.Session
	// config           *AppConfig
	config   ConfigMap
	plugins  PluginRegistry
	services []ServiceProvider
	//routeMiddlewares map[string]Middleware
	//httpMiddlewares  []func(http.Handler) http.Handler
	//routes           []*Route
	hooks          *AppHooks
	router         *Router
	db             *db.DB
	dbFunc         func(c context.Context, config *db.DBConfig) (*db.DB, error)
	i              *inertia.Inertia
	routeRegistrar RouteRegistrarFunc
	currentGroup   *RouteGroup
	fs             fsys.FS
}

type Options struct {
	container.Container
	*session.Session
	Config           ConfigMap
	Plugins          map[PluginID]Plugin
	Providers        []ServiceProvider
	routeMiddlewares map[string]Middleware
	Hooks            *AppHooks
	inertia          *inertia.Inertia
	fs               fsys.FS
}

type OptFunc func(opts *Options)

func (app *App) Plugin(namespace string) Plugin {
	return app.plugins.Find(namespace)
}

func (app *App) RegisterRoutes(fn RouteRegistrarFunc) {
	app.router.routeRegistrar = fn
}

func (app *App) Router() *Router {
	return app.router
}

func (app *App) Session() *session.Session {
	return app.session
}

func (app *App) Inertia() *inertia.Inertia {
	return app.i
}

func (app *App) DB() *db.DB {
	return app.db
}

func (app *App) DbFunc(c context.Context, config *db.DBConfig) (*db.DB, error) {
	return app.dbFunc(c, config)
}

func (app *App) FS() fsys.FS {
	return app.fs
}

func getDefaultConfig() ConfigMap {
	return ConfMap()
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
		nil,
		nil,
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

func WithInertia(i *inertia.Inertia) OptFunc {
	if i == nil {
		i = initInertia()
	}
	return func(opts *Options) {
		opts.inertia = i
	}
}

func WithFS(fs fsys.FS) OptFunc {
	if fs == nil {
		fs = fsys.NewLocalStorage("./storage")
	}
	return func(opts *Options) {
		opts.fs = fs
	}
}

//func WithRouter(router HTTPRouter) OptFunc {
//	return func(opts *Options) {
//		opts.HTTPRouter = router
//	}
//}

func WithContainer(container container.Container) OptFunc {
	return func(opts *Options) {
		opts.Container = container
	}
}

func WithSession(sm *session.Session) OptFunc {
	return func(opts *Options) {
		opts.Session = sm
	}
}

func NewApp(optFuncs ...OptFunc) *App {
	opts := defaultOptions()
	hr := NewRouter(NewChiRouter())

	for _, optFunc := range optFuncs {
		optFunc(opts)
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
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					panic(err)
				}
				slog.Info("Copied template %s to %s\n", name, filePath)
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

	hr.setRouteMiddleware(routeMiddlewares)

	app := &App{
		// opts.Container,
		isContextReady: false,
		mu:             sync.Mutex{},
		config:         opts.Config,
		plugins:        opts.Plugins,
		services:       opts.Providers,
		hooks:          opts.Hooks,
		router:         hr,
		i:              opts.inertia,
		fs:             opts.fs,
	}
	return app
}

// func (app *App) Container() container.Container {
// 	return app.container
// }

func (app *App) registerServices(services []ServiceProvider) {
	for _, svc := range services {
		extendsBase := false
		if reflect.TypeOf(svc).Kind() != reflect.Ptr {
			panic("Service must be a pointer")
		}
		if reflect.TypeOf(svc).Elem().Kind() != reflect.Struct {
			panic("Service must be a struct")
		}

		// Iterate over all the fields of the struct and see if it extends BaseServiceProvider
		for i := 0; i < reflect.TypeOf(svc).Elem().NumField(); i++ {
			if reflect.TypeOf(svc).Elem().Field(i).Type == reflect.TypeOf(BaseServiceProvider{}) {
				extendsBase = true
				break
			}
		}

		if !extendsBase {
			panic("Service must extend BaseServiceProvider")
		}

		// Check if service implements ServiceProvider interface, not necessary if type hinted
		if reflect.TypeOf(svc).Implements(reflect.TypeOf((*ServiceProvider)(nil)).Elem()) {
			slog.Info("Registering service: " + reflect.TypeOf(svc).Elem().Name())
			svc.Register(app)
			app.services = append(app.services, svc)
		} else {
			panic("Service must implement ServiceProvider interface")
		}
	}

	for _, service := range services {
		if reflect.TypeOf(service).Implements(reflect.TypeOf((*ServiceProvider)(nil)).Elem()) {
			service.Boot()
		}
	}
}

func (app *App) Run() {
	app.registerServices([]ServiceProvider{
		&DatabaseServiceProvider{},
		&SessionServiceProvider{},
		&AuthServiceProvider{},
	})

	app.router.registerMiddlewares(app)

	app.router.registerRoutes(app)

	app.router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	app.router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	for _, plugin := range app.plugins {
		if err := plugin.Boot(app); err != nil {
			panic(err)
		}
	}

	slog.Info(fmt.Sprintf("%s is running on port %d...", app.config.get("app.name"), app.config.get("app.port")))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", app.config.get("app.port")), app.session.LoadAndSave(app.router)); err != nil {
		panic(err)
	}
}

func (app *App) HandleSignals() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	sig := <-signalChannel
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		app.Shutdown()
		os.Exit(0)
	}
}

func (app *App) Shutdown() {
	log.Println("Shutting down application...")
	sessName := app.db.Name()
	err := app.db.Close()
	if err != nil {
		log.Fatal("Error closing database connection:", err)
	}
	log.Println("Database connection", sessName, "closed.")
}
