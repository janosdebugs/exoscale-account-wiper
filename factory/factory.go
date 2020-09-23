package factory

import (
	"github.com/janoszen/exoscale-account-wiper/aa"
	"github.com/janoszen/exoscale-account-wiper/dns"
	"github.com/janoszen/exoscale-account-wiper/eips"
	"github.com/janoszen/exoscale-account-wiper/iam"
	"github.com/janoszen/exoscale-account-wiper/instances"
	"github.com/janoszen/exoscale-account-wiper/nlbs"
	"github.com/janoszen/exoscale-account-wiper/pluginregistry"
	"github.com/janoszen/exoscale-account-wiper/pools"
	"github.com/janoszen/exoscale-account-wiper/privnets"
	"github.com/janoszen/exoscale-account-wiper/sg"
	"github.com/janoszen/exoscale-account-wiper/sos"
	"github.com/janoszen/exoscale-account-wiper/sshkeys"
	"github.com/janoszen/exoscale-account-wiper/templates"
)

func CreateRegistry() *pluginregistry.PluginRegistry {
	r := pluginregistry.New()
	r.Register(eips.New())
	r.Register(nlbs.New())
	r.Register(pools.New())
	r.Register(sg.New())
	r.Register(instances.New())
	r.Register(templates.New())
	r.Register(aa.New())
	r.Register(sshkeys.New())
	r.Register(privnets.New())
	r.Register(sos.New())
	r.Register(dns.New())
	r.Register(iam.New())
	return r
}
