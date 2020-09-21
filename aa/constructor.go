package aa

import "github.com/janoszen/exoscale-account-wiper/plugin"

func New() plugin.DeletePlugin {
	return &Plugin{}
}
