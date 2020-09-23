package instances_test

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/instances"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingInstances(t *testing.T) {
	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()
	clientFactory := plugin.NewClientFactory(tf.ExoscaleKey, tf.ExoscaleSecret)
	v1Client := clientFactory.GetExoscaleClient()
	instancePrototype := &egoscale.VirtualMachine{}
	sgs, err := v1Client.ListWithContext(context.Background(), instancePrototype)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed to list instances (%v)", err))
	}
	assert.Equal(t, 1, len(sgs), fmt.Sprintf("invalid number of instances returned (%d)", len(sgs)))

	i := instances.New()
	err = i.Run(clientFactory, context.Background())
	if err != nil {
		t.Fail()
	}

	sgs, err = v1Client.ListWithContext(context.Background(), instancePrototype)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed to list instances (%v)", err))
	}
	assert.Equal(t, 0, len(sgs), fmt.Sprintf("invalid number of instances returned (%d)", len(sgs)))
}
