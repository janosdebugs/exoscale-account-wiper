package eips

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"log"
)

type Plugin struct {
}

func (p *Plugin) GetKey() string {
	return "eips"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("EIP deletion has no options")
}

func (p *Plugin) Run(clientFactory *plugin.ClientFactory, ctx context.Context) error {
	log.Printf("deleting EIP's...")

	client := clientFactory.GetExoscaleClient()
	req := egoscale.IPAddress{}
	ips, err := client.ListWithContext(ctx, &req)
	if err != nil {
		return err
	}

	for _, ip := range ips {
		eip := ip.(*egoscale.IPAddress)
		if !eip.IsElastic {
			continue
		}
		log.Printf("deleting EIP %s...", eip.ID)
		err := eip.Delete(ctx, client)
		if err != nil {
			log.Printf("failed to delete EIP %s (%v)", eip.ID, err)
		} else {
			log.Printf("deleted EIP %s.", eip.ID)
		}
	}

	log.Printf("deleted EIP's.")
	return nil
}
