package templates

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
	return "templates"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("template deletion has no options")
}

func (p *Plugin) Run(clientFactory *plugin.ClientFactory, ctx context.Context) error {
	log.Printf("deleting templates...")

	client := clientFactory.GetExoscaleClient()
	resp, err := client.RequestWithContext(ctx, egoscale.ListZones{})
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	poolBlocker := make(chan bool, 10)
	for _, z := range resp.(*egoscale.ListZonesResponse).Zone {
		req := &egoscale.ListTemplates{
			TemplateFilter: "self",
			ZoneID:         z.ID,
			Keyword:        "",
		}
		client.PaginateWithContext(ctx, req, func(i interface{}, e error) bool {
			template := i.(*egoscale.Template)
			wg.Add(1)
			go func() {
				defer wg.Done()
				poolBlocker <- true

				log.Printf("deleting template %s...", template.ID)
				cmd := egoscale.DeleteTemplate{
					ID: template.ID,
				}
				err := client.BooleanRequestWithContext(ctx, cmd)
				if err != nil {
					log.Printf("error while deleting tempalte %s (%v)", template.ID, err)
				} else {
					log.Printf("deleted template %s...", template.ID)
				}

				<-poolBlocker
			}()
			return true
		})
	}
	log.Printf("deleted templates.")
	return nil
}
