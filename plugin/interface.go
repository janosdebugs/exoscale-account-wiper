package plugin

import (
	"context"
)

type DeletePlugin interface {
	// Return a unique key for this plugin
	GetKey() string
	// Return a parameter list with descriptions.
	GetParameters() map[string]string
	// Set a parameter. The plugin may return an error if the configuration value is not valid or no such configuration
	// option exists. Parameter names will be passed lowercase, separated by a dash (-) even when the input was sent
	// with an underscore (-)
	SetParameter(name string, value string) error
	// Run the deletion process. Will only be called if a deletion of that particular resource is requested. Should
	// handle failures to delete the resource.
	Run(client *ClientFactory, ctx context.Context) error
}
