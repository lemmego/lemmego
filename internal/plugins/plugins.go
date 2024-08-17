package plugins

import (
	"github.com/lemmego/lemmego/api/app"
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

// Load plugins
func Load() app.PluginRegistry {
	authPlugin := auth.New()
	registry.Add(authPlugin)
	return registry
}
