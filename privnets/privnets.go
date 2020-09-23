package privnets

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"log"
	"sync"
)

type Plugin struct {
}

func (p *Plugin) GetKey() string {
	return "privnets"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("privnet deletion has no options")
}

func (p *Plugin) Run(clientFactory *plugin.ClientFactory, ctx context.Context) error {
	log.Printf("deleting private networks...")

	client := clientFactory.GetExoscaleClient()
	var wg sync.WaitGroup
	poolBlocker := make(chan bool, 10)

	zones, err := client.ListWithContext(ctx, &egoscale.Zone{})
	if err != nil {
		return err
	}

	for _, z := range zones {
		req := egoscale.Network{
			ZoneID:          z.(*egoscale.Zone).ID,
			Type:            "Isolated",
			CanUseForDeploy: true,
		}

		zoneName := z.(*egoscale.Zone).Name
		privnets, err := client.ListWithContext(ctx, &req)
		if err != nil {
			log.Printf("failed to list private networks in zone %s (%v)", zoneName, err)
			continue
		}

		for _, p := range privnets {
			privnet := p.(*egoscale.Network)
			wg.Add(1)
			go func() {
				defer wg.Done()
				poolBlocker <- true
				defer func() { <-poolBlocker }()

				log.Printf("deleting private network %s in zone %s...", privnet.ID, zoneName)
				addrReq := &egoscale.DeleteNetwork{
					ID: privnet.ID,
				}
				err := client.BooleanRequestWithContext(ctx, addrReq)
				if err != nil {
					log.Printf("failed to delete private network %s in zone %s (%v)", privnet.ID, zoneName, err)
				} else {
					log.Printf("deleted private network %s in zone %s.", privnet.ID, zoneName)
				}
			}()
		}
	}

	wg.Wait()
	log.Printf("deleted private networks.")
	return nil
}
