package pluginregistry

import "github.com/janoszen/exoscale-account-wiper/plugin"

func New() *PluginRegistry {
	return &PluginRegistry{
		plugins:      []plugin.DeletePlugin{},
		pluginsByKey: make(map[string]plugin.DeletePlugin),
		enabledPlugins: make(map[string]bool),
	}
}
