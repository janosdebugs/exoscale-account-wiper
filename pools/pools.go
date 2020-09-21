package pools

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"log"
	"sync"
	"time"
)

type Plugin struct {
	logger log.Logger
}

func (p *Plugin) GetKey() string {
	return "pools"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("instance pool deletion has no options")
}

func (p *Plugin) Run(client *egoscale.Client, ctx context.Context) error {
	log.Printf("deleting instance pools...")

	var wg sync.WaitGroup
	poolBlocker := make(chan bool, 10)

	select {
	case <-ctx.Done():
		log.Printf("aborting...")
		return nil
	default:
	}

	resp, err := client.RequestWithContext(ctx, egoscale.ListZones{})
	if err != nil {
		return err
	}
	for _, z := range resp.(*egoscale.ListZonesResponse).Zone {
		select {
		case <-ctx.Done():
			break
		default:
		}

		resp, err := client.RequestWithContext(ctx, egoscale.ListInstancePools{ZoneID: z.ID})
		if err != nil {
			return err
		}
		for _, i := range resp.(*egoscale.ListInstancePoolsResponse).InstancePools {
			select {
			case <-ctx.Done():
				break
			default:
			}

			wg.Add(1)
			instancePoolId := i.ID
			zoneId := z.ID
			currentState := i.State
			go func() {
				defer wg.Done()
				poolBlocker <- true
				defer func() { <-poolBlocker }()
				log.Printf("deleting instance pool %s...", instancePoolId)
				var err error = nil
				if currentState != egoscale.InstancePoolDestroying {
					request := egoscale.DestroyInstancePool{
						ID:     instancePoolId,
						ZoneID: zoneId,
					}
					err = client.BooleanRequestWithContext(ctx, request)
				} else {
					log.Printf("instance pool %s is already being destroyed...", instancePoolId)
				}

				if err != nil {
					log.Printf("error deleting instance pool %s (%v)", instancePoolId, err)
				} else {
					for {
						log.Printf("waiting for complete removal of instance pool %s...", instancePoolId)
						getRequest := egoscale.GetInstancePool{
							ID:     instancePoolId,
							ZoneID: zoneId,
						}
						if _, err := client.RequestWithContext(ctx, getRequest); err != nil {
							//Wait for the instance pool to be completely destroyed
							log.Printf("deleted instance pool %s", instancePoolId)
							break
						}
						time.Sleep(time.Second * 10)
					}
				}
			}()
		}
	}
	wg.Wait()
	log.Printf("deleted instance pools.")

	return nil
}
