package pools_test

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/pools"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingInstancePool(t *testing.T) {
	t.Skip("Test skipped because the teardown fails due to the following bug: https://github.com/exoscale/terraform-provider-exoscale/issues/72")
	return

	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()

	client := egoscale.NewClient("https://api.exoscale.ch/v1", tf.ExoscaleKey, tf.ExoscaleSecret)

	resp, err := client.Request(egoscale.ListZones{})
	if err != nil {
		t.Fail()
		return
	}
	zones := resp.(*egoscale.ListZonesResponse).Zone
	instancePoolCount := 0
	for _, z := range zones {
		resp, err := client.Request(egoscale.ListInstancePools{ZoneID: z.ID})
		if err != nil {
			t.Fail()
			return
		}
		instancePoolCount += len(resp.(*egoscale.ListInstancePoolsResponse).InstancePools)
	}
	assert.Equal(t, 1, instancePoolCount, fmt.Sprintf("invalid number of instance pools returned (%d)", instancePoolCount))

	i := pools.New()
	err = i.Run(client, context.Background())
	if err != nil {
		t.Fail()
	}

	instancePoolCount = 0
	for _, z := range zones {
		resp, err := client.Request(egoscale.ListInstancePools{ZoneID: z.ID})
		if err != nil {
			t.Fail()
			return
		}
		instancePoolCount += len(resp.(*egoscale.ListInstancePoolsResponse).InstancePools)
	}
	assert.Equal(t, 0, instancePoolCount, fmt.Sprintf("invalid number of instance pools returned (%d)", instancePoolCount))
}
