package nlbs

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/exoscale/egoscale"
	apiv2 "github.com/exoscale/egoscale/api/v2"
)

type Plugin struct {
	logger log.Logger
}

func (p *Plugin) GetKey() string {
	return "nlbs"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("NLB deletion has no options")
}

func (p *Plugin) Run(client *egoscale.Client, ctx context.Context) error {
	log.Printf("deleting NLB's...")

	resp, err := client.RequestWithContext(ctx, egoscale.ListZones{})
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	poolBlocker := make(chan bool, 10)
	for _, z := range resp.(*egoscale.ListZonesResponse).Zone {
		select {
		case <-ctx.Done():
			break
		default:
		}

		zoneName := z.Name
		v2Context := apiv2.WithEndpoint(ctx, apiv2.NewReqEndpoint("", z.Name))
		nlbs, err := client.ListNetworkLoadBalancers(v2Context, z.Name)
		if err != nil {
			log.Printf("failed to list NLB's in zone %s (%v)", z.Name, err)
			continue
		}
		for _, nlb := range nlbs {
			select {
			case <-ctx.Done():
				break
			default:
			}

			nlbId := nlb.ID
			nlbState := nlb.State
			go func() {
				wg.Add(1)
				defer wg.Done()
				poolBlocker <- true
				defer func() { <-poolBlocker }()

				if nlbState == "Deleting" {
					log.Printf("NLB %s in zone %s is already being deleted", nlbId, zoneName)
					for {
						log.Printf("waiting for complete removal of NLB %s in zone %s", nlbId, zoneName)
						_, err := client.GetNetworkLoadBalancer(v2Context, zoneName, nlbId)
						if err != nil {
							break
						}
						time.Sleep(time.Second * 10)
					}
				} else {
					err := client.DeleteNetworkLoadBalancer(v2Context, zoneName, nlbId)
					if err != nil {
						log.Printf("failed to delete NLB %s in zone %s (%v)", nlbId, zoneName, err)
					} else {
						log.Printf("deleted NLB %s in zone %s", nlbId, zoneName)
					}
				}
			}()
		}
	}
	wg.Wait()
	log.Printf("deleted NLB's.")

	return nil
}
