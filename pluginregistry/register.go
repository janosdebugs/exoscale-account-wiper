package pluginregistry

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"strings"
)

func (r * PluginRegistry) Register(plugin plugin.DeletePlugin) {
	r.plugins = append(r.plugins, plugin)
	r.pluginsByKey[plugin.GetKey()] = plugin
	r.enabledPlugins[plugin.GetKey()] = true
}

func (r * PluginRegistry) SetConfiguration(config map[string]string, ignoreErrors bool) error {
	for key, value := range config {
		var keyParts []string
		if strings.Contains(key, "_") {
			//Assume underscore separators
			keyParts = strings.Split(key, "_")
		} else if strings.Contains(key, "-") {
			//Assume dash separators
			keyParts = strings.Split(key, "-")
		} else if ignoreErrors {
			continue
		} else {
			return fmt.Errorf("invalid configuration option %s", key)
		}
		module := strings.ToLower(keyParts[0])
		parameterName := strings.ToLower(strings.Join(keyParts[1:], "-"))
		if p, ok := r.pluginsByKey[module]; ok {
			err := p.SetParameter(parameterName, value)
			if err != nil && !ignoreErrors {
				return err
			}
		} else if ignoreErrors {
			continue;
		} else {
			return fmt.Errorf("invalid configuration option %s, no module named %s", key, module)
		}
	}
	return nil
}

func (r *PluginRegistry) GetPlugins() map[string]plugin.DeletePlugin {
	return r.pluginsByKey
}

func (r *PluginRegistry) EnablePlugin(plugin string) error {
	if _, ok := r.enabledPlugins[plugin]; !ok {
		return fmt.Errorf("no such plugin: %s", plugin)
	}
	r.enabledPlugins[plugin] = true
	return nil
}

func (r *PluginRegistry) DisablePlugin(plugin string) error {
	if _, ok := r.enabledPlugins[plugin]; !ok {
		return fmt.Errorf("no such plugin: %s", plugin)
	}
	r.enabledPlugins[plugin] = false
	return nil
}

func (r * PluginRegistry) Run(v1Client * egoscale.Client, ctx context.Context) error {
	for _, p := range r.plugins {
		if r.enabledPlugins[p.GetKey()] {
			err := p.Run(v1Client, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
