package iam

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"log"
)

type Plugin struct {
	excludeSelf bool
}

func (p *Plugin) GetKey() string {
	return "iam"
}

func (p *Plugin) GetParameters() map[string]string {
	return map[string]string{
		"exclude-self": "Exclude current API key from deletion",
	}
}

func (p *Plugin) SetParameter(name string, value string) error {
	if name != "exclude-self" {
		return fmt.Errorf("invalid option: %s", name)
	}
	if value == "true" || value == "1" {
		p.excludeSelf = true
	} else {
		p.excludeSelf = false
	}
	return nil
}

func (p *Plugin) Run(clientFactory *plugin.ClientFactory, ctx context.Context) error {
	log.Printf("deleting IAM keys...")

	client := clientFactory.GetExoscaleClient()
	resp, err := client.RequestWithContext(ctx, &egoscale.ListAPIKeys{})
	if err != nil {
		return err
	}

	r := resp.(*egoscale.ListAPIKeysResponse)
	for _, apiKey := range r.APIKeys {
		if p.excludeSelf && apiKey.Key == client.APIKey {
			log.Printf("skipping deletion of current API key %s.", apiKey.Name)
			continue
		}
		log.Printf("deleting API key %s...", apiKey.Name)

		response, err := client.RequestWithContext(ctx, &egoscale.RevokeAPIKey{Key: apiKey.Key})
		if err != nil {
			log.Printf("failed to revoke API key %s (%v)", apiKey.Name, err)
		} else {
			revokeResponse := response.(*egoscale.RevokeAPIKeyResponse)
			if revokeResponse.Success {
				log.Printf("deleted API key %s.", apiKey.Name)
			} else {
				log.Printf("failed to deleted API key %s.", apiKey.Name)
			}
		}
	}

	log.Printf("deleted IAM keys.")
	return nil
}
