package privnets_test

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/privnets"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingPrivnets(t *testing.T) {
	t.Skip("Skipping test because of Exoscale Terraform provider bug https://github.com/exoscale/terraform-provider-exoscale/issues/73")
	return

	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()

	v1Client := egoscale.NewClient("https://api.exoscale.ch/v1", tf.ExoscaleKey, tf.ExoscaleSecret)
	zones, err := v1Client.ListWithContext(context.Background(), &egoscale.Zone{})
	if err != nil {
		assert.Fail(t, "error while listing zones (%v)", err)
	}
	privNetCount := 0
	for _, z := range zones {
		req := egoscale.Network{
			ZoneID:          z.(*egoscale.Zone).ID,
			Type:            "Isolated",
			CanUseForDeploy: true,
		}

		zoneName := z.(*egoscale.Zone).Name
		privNets, err := v1Client.ListWithContext(context.Background(), &req)
		if err != nil {
			assert.Fail(t, "error while listing privnets in zone %s (%v)", zoneName, err)
		} else {
			privNetCount += len(privNets)
		}
	}
	assert.Equal(t, 1, privNetCount, fmt.Sprintf("invalid number of instances returned (%d)", privNetCount))

	i := privnets.New()
	err = i.Run(v1Client, context.Background())
	if err != nil {
		t.Fail()
	}

	zones, err = v1Client.ListWithContext(context.Background(), &egoscale.Zone{})
	if err != nil {
		assert.Fail(t, "error while listing zones (%v)", err)
	}
	privNetCount = 0
	for _, z := range zones {
		req := egoscale.Network{
			ZoneID:          z.(*egoscale.Zone).ID,
			Type:            "Isolated",
			CanUseForDeploy: true,
		}

		zoneName := z.(*egoscale.Zone).Name
		privNets, err := v1Client.ListWithContext(context.Background(), &req)
		if err != nil {
			assert.Fail(t, "error while listing privnets in zone %s (%v)", zoneName, err)
			privNetCount += len(privNets)
		}
	}
	assert.Equal(t, 0, privNetCount, fmt.Sprintf("invalid number of instances returned (%d)", privNetCount))

}
