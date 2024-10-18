package plugins

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/lemmego/internal/plugins/auth"
)

var registry = app.PluginRegistry{}

// Add a plugin to the registry
func Add(plugin app.Plugin) {
	registry.Add(plugin)
}

// Get a plugin from the registry
func Get[T app.Plugin](plugin T) T {
	return registry.Get(plugin).(T)
}

// AddByPkgName adds a plugin to the registry by the package name
func AddByPkgName(id string, plugin app.Plugin) {
	registry.Add(plugin)
}

// GetByPkgName gets a plugin to the registry by the package name
func GetByPkgName(id string) app.Plugin {
	return registry[app.PluginID(id)]
}

func init() {
	Add(auth.New())
}

// Load plugins
func Load() app.PluginRegistry {
	return registry
}
