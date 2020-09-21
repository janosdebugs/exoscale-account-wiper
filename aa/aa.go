package aa

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"log"
)

type Plugin struct {
}

func (p *Plugin) GetKey() string {
	return "aa"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("anti-affinity group deletion has no options")
}

func (p *Plugin) Run(client *egoscale.Client, ctx context.Context) error {
	log.Printf("deleting anti-affinity groups...")

	resp, err := client.RequestWithContext(ctx, &egoscale.ListAffinityGroups{})
	if err != nil {
		return fmt.Errorf("failed to list affinity groups (%v)", err)
	}

	for _, ag := range resp.(*egoscale.ListAffinityGroupsResponse).AffinityGroup {
		log.Printf("deleting anti-affinity group %s...", ag.ID)
		err := ag.Delete(ctx, client)
		if err != nil {
			log.Printf("failed to delete affinity group %s (%v)", ag.ID, err)
		} else {
			log.Printf("deleted anti-affinity group %s.", ag.ID)
		}
	}

	log.Printf("deleted anti-affinity groups.")
	return nil
}
