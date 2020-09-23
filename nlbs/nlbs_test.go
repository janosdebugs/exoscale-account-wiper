package nlbs_test

import (
	"context"
	apiv2 "github.com/exoscale/egoscale/api/v2"
	"github.com/janoszen/exoscale-account-wiper/nlbs"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingNetworkLoadBalancer(t *testing.T) {
	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()
	clientFactory := plugin.NewClientFactory(tf.ExoscaleKey, tf.ExoscaleSecret)
	client := clientFactory.GetExoscaleClient()

	v2Context := apiv2.WithEndpoint(context.Background(), apiv2.NewReqEndpoint("", "at-vie-1"))
	nlbList, err := client.ListNetworkLoadBalancers(v2Context, "at-vie-1")
	if err != nil {
		assert.Fail(t, "failed to list NLB's", err)
	}
	assert.Equal(t, 1, len(nlbList), "invalid NLB count after setup: %d", len(nlbList))

	i := nlbs.New()
	err = i.Run(clientFactory, context.Background())
	if err != nil {
		t.Fail()
	}

	nlbList, err = client.ListNetworkLoadBalancers(v2Context, "at-vie-1")
	if err != nil {
		assert.Fail(t, "failed to list NLB's", err)
	}
	assert.Equal(t, 0, len(nlbList), "invalid NLB count after run: %d", len(nlbList))
}
