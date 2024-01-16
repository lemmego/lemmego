package plugins

import (
	"pressebo/framework"
	"pressebo/plugins/auth"
)

// Load plugins
func Load() framework.PluginRegistry {
	registry := framework.PluginRegistry{}
	authPlugin := auth.New()

	registry.Add(authPlugin)
	return registry
}
