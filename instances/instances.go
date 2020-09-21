package instances

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"log"
	"sync"
	"time"
)

type Plugin struct {
}

func (p *Plugin) GetKey() string {
	return "instances"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("instance deletion has no options")
}

func (p *Plugin) Run(client *egoscale.Client, ctx context.Context) error {
	log.Printf("deleting instances...")
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
		wg.Add(1)

		go func() {
			poolBlocker <- true
			if vm.State == "Destroying" {
				log.Printf("instance %s is already being destroyed.\n", vm.ID)
			} else {
				log.Printf("deleting instance %s...\n", vm.ID)
				err := vm.Delete(ctx, client)
				if err != nil {
					log.Printf("could not delete instance %s (%v)\n", vm.ID, err)
					return
				}
			}

			for {
				log.Printf("waiting for instance %s to be destroyed...\n", vm.ID)
				_, err := client.Get(vm)
				if err != nil {
					break
				}
				time.Sleep(time.Second * 10)
			}
			log.Printf("deleted instance %s\n", vm.ID)
			<-poolBlocker
			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("deleted instances.")
	return nil
}
