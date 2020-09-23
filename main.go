package main

import (
	"context"
	"fmt"
	"github.com/janoszen/exoscale-account-wiper/aa"
	"github.com/janoszen/exoscale-account-wiper/dns"
	"github.com/janoszen/exoscale-account-wiper/eips"
	"github.com/janoszen/exoscale-account-wiper/instances"
	"github.com/janoszen/exoscale-account-wiper/nlbs"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"github.com/janoszen/exoscale-account-wiper/pluginregistry"
	"github.com/janoszen/exoscale-account-wiper/pools"
	"github.com/janoszen/exoscale-account-wiper/privnets"
	"github.com/janoszen/exoscale-account-wiper/sg"
	"github.com/janoszen/exoscale-account-wiper/sos"
	"github.com/janoszen/exoscale-account-wiper/sshkeys"
	"github.com/janoszen/exoscale-account-wiper/templates"
	"log"
	"os"
	"strings"
)

func createRegistry() *pluginregistry.PluginRegistry {
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
	return r
}

func usage(registry *pluginregistry.PluginRegistry) {
	fmt.Printf("Usage: exoscale-account-wiper [OPTIONS]\n")
	fmt.Printf("\n")
	fmt.Printf("Arguments can be passed via a command line option or via an environment variable.\n")
	fmt.Printf("\n")
	fmt.Printf("API credentials:\n")
	fmt.Printf("  --apikey API-KEY-HERE or APIKEY=API-KEY-HERE:  Pass Exoscale API key\n")
	fmt.Printf("  --apisecret API-SECRET-HERE or APISECRET=API-SECRET-HERE:  Pass Exoscale API secret\n")
	fmt.Printf("\n")
	fmt.Printf("Enable/disable plugins:\n")
	fmt.Printf("  --[no]delete or DELETE=0|1:  Enable or disable deletion modules by default\n")
	for _, p := range registry.GetPlugins() {
		fmt.Printf("  --[no]%s or %s=0|1:  Enable or disable %s deletion\n", p.GetKey(), strings.Replace(strings.ToUpper(p.GetKey()), "-", "_", -1), p.GetKey())
	}
	for _, p := range registry.GetPlugins() {
		parameters := p.GetParameters()
		if len(parameters) > 0 {
			fmt.Printf("\n")
			fmt.Printf("%s OPTIONS\n", strings.ToUpper(p.GetKey()))
			fmt.Printf("\n")
			for param, description := range parameters {
				fmt.Printf("  --%s-%s or %s_%s: %s\n", p.GetKey(), param, strings.ToUpper(p.GetKey()), strings.Replace(strings.ToUpper(param), "-", "_", -1), description)
			}
		}
	}
}

func main() {
	registry := createRegistry()

	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		usage(registry)
		return
	}

	envOptions := map[string]string{}
	environment := os.Environ()
	for _, line := range environment {
		parts := strings.SplitN(line, "=", 2)
		key := parts[0]
		value := parts[1]
		envOptions[key] = value
	}
	err := registry.SetConfiguration(envOptions, true)
	if err != nil {
		log.Fatal(err)
	}

	plugins := registry.GetPlugins()
	i := 1
	enabledPlugins := map[string]*bool{}
	for p := range plugins {
		enabledPlugins[p] = nil
	}
	defaultEnabled := true
	t := true
	f := false
	options := map[string]string{}
	for {
		if len(os.Args) < i+1 {
			break
		}
		item := os.Args[i]
		if item[:2] == "--" {
			option := strings.ToLower(item[2:])
			if _, ok := enabledPlugins[option]; ok {
				enabledPlugins[option] = &t
			} else if _, ok := enabledPlugins[option[2:]]; ok && option[:2] == "no" {
				enabledPlugins[option] = &f
			} else if option == "delete" {
				defaultEnabled = true
			} else if option == "nodelete" {
				defaultEnabled = false
			} else {
				nextItem := os.Args[i+1]
				if nextItem[:2] == "--" {
					options[option] = "0"
				} else {
					options[option] = nextItem
					i++
				}
			}
		} else {
			fmt.Printf("invalid parameter: %s\n", item)
			usage(registry)
			os.Exit(1)
		}
		i++
	}
	for enabledPlugin := range enabledPlugins {
		if (enabledPlugins[enabledPlugin] != nil && *enabledPlugins[enabledPlugin]) ||
			(enabledPlugins[enabledPlugin] == nil && defaultEnabled) {
			err := registry.EnablePlugin(enabledPlugin)
			if err != nil {
				fmt.Printf("Error: unable to activate plugin %s (%v)", enabledPlugin, err)
				usage(registry)
				os.Exit(1)
			}
		} else {
			err := registry.DisablePlugin(enabledPlugin)
			if err != nil {
				fmt.Printf("Error: unable to deactivate plugin %s (%v)", enabledPlugin, err)
				usage(registry)
				os.Exit(1)
			}
		}
	}

	var exoscaleApiKey string
	var exoscaleApiSecret string
	var ok bool
	if exoscaleApiKey, ok = options["apikey"]; ok {
		delete(options, "apikey")
	} else if val, ok := envOptions["APIKEY"]; ok {
		exoscaleApiKey = val
	} else {
		fmt.Printf("Error: no API key provided\n")
		usage(registry)
		os.Exit(1)

	}
	if exoscaleApiSecret, ok = options["apisecret"]; ok {
		delete(options, "apisecret")
	} else if val, ok := envOptions["APISECRET"]; ok {
		exoscaleApiSecret = val
	} else {
		fmt.Printf("Error: no API key provided\n")
		usage(registry)
		os.Exit(1)
	}

	err = registry.SetConfiguration(options, false)
	if err != nil {
		log.Fatal(err)
	}

	clientFactory := plugin.NewClientFactory(exoscaleApiKey, exoscaleApiSecret)

	ctx := context.Background()
	err = registry.Run(clientFactory, ctx)
	if err != nil {
		log.Fatal(err)
	}
}
