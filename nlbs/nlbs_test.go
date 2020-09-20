package nlbs_test

import (
	"context"
	"github.com/exoscale/egoscale"
	apiv2 "github.com/exoscale/egoscale/api/v2"
	"github.com/janoszen/exoscale-account-wiper/nlbs"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingNetworkLoadBalancer(t *testing.T) {
	t.Skip("Skipping test because of Exoscale Terraform provider bug https://github.com/exoscale/terraform-provider-exoscale/issues/73")
	return

	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()

	client := egoscale.NewClient("https://api.exoscale.ch/v1", tf.ExoscaleKey, tf.ExoscaleSecret)

	v2Context := apiv2.WithEndpoint(context.Background(), apiv2.NewReqEndpoint("", "at-vie-1"))
	nlbList, err := client.ListNetworkLoadBalancers(v2Context, "at-vie-1")
	if err != nil {
		assert.Fail(t, "failed to list NLB's", err)
	}
	assert.Equal(t, 1, len(nlbList), "invalid NLB count after setup: %d", len(nlbList))

	i := nlbs.New()
	err = i.Run(client, context.Background())
	if err != nil {
		t.Fail()
	}

	nlbList, err = client.ListNetworkLoadBalancers(v2Context, "at-vie-1")
	if err != nil {
		assert.Fail(t, "failed to list NLB's", err)
	}
	assert.Equal(t, 0, len(nlbList), "invalid NLB count after run: %d", len(nlbList))
}
