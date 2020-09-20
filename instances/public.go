package instances

import (
	"context"
	"github.com/exoscale/egoscale"
	"log"
	"sync"
)

func (p * Plugin) GetKey() string {
	return "instances"
}

func (p *  Plugin) Run(client *egoscale.Client, ctx context.Context) error {
	log.Printf("deleting virtual machines...")
	vm := &egoscale.VirtualMachine{}
	vms, err := client.ListWithContext(ctx, vm)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	poolBlocker := make(chan bool, 10)

	for _, key := range vms {
		select {
		case <-ctx.Done():
			break
		default:
		}
		vm := key.(*egoscale.VirtualMachine)

		go func() {
			wg.Add(1)
			poolBlocker <- true
			log.Printf("deleting instance %s...\n", vm.ID)
			err := vm.Delete(ctx, client)
			if err != nil {
				log.Printf("could not delete instance %s (%v)\n", vm.ID, err)
			} else {
				log.Printf("deleted instance %s\n", vm.ID)
			}
			<- poolBlocker
		}()
	}
	wg.Wait()
	log.Printf("deleted virtual machines.")
	return nil
}
