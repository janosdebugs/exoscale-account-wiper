package sshkeys

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
	return "sg"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("security group deletion has no options")
}

func (p *Plugin) Run(clientFactory *plugin.ClientFactory, ctx context.Context) error {
	log.Printf("deleting SSH keys...")

	client := clientFactory.GetExoscaleClient()
	sshKeys, err := client.ListWithContext(ctx, &egoscale.SSHKeyPair{})
	if err != nil {
		return err
	}
	for _, key := range sshKeys {
		sshKey := key.(*egoscale.SSHKeyPair)
		log.Printf("deleting SSH key %s...", sshKey.Name)
		err := sshKey.Delete(ctx, client)
		if err != nil {
			log.Printf("failed to delete SSH key %s (%v)", sshKey.Name, err)
		} else {
			log.Printf("deleted SSH key %s.", sshKey.Name)
		}
	}

	log.Printf("deleted SSH keys.")
	return nil
}
