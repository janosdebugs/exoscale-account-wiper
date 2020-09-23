package dns

import (
	"context"
	"fmt"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"log"
)

type Plugin struct {
}

func (p *Plugin) GetKey() string {
	return "dns"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("DNS zone deletion has no options")
}

func (p *Plugin) Run(clientFactory *plugin.ClientFactory, ctx context.Context) error {
	log.Printf("deleting DNS zones...")

	client := clientFactory.GetDnsClient()

	domains, err := client.GetDomains(ctx)
	if err != nil {
		return err
	}
	for _, domain := range domains {
		log.Printf("deleting domain %s...", domain.Name)
		err = client.DeleteDomain(ctx, domain.Name)
		if err != nil {
			log.Printf("failed to delete domain %s (%v)", domain.Name, err)
		} else {
			log.Printf("deleted domain %s.", domain.Name)
		}
	}

	log.Printf("deleted DNS zones.")
	return nil
}
