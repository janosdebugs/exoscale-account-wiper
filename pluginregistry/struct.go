package pluginregistry

import "github.com/janoszen/exoscale-account-wiper/plugin"

type PluginRegistry struct {
	plugins []plugin.DeletePlugin
	pluginsByKey map[string]plugin.DeletePlugin
	enabledPlugins map[string]bool
}
