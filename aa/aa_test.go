package aa_test

import (
	"context"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/aa"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingAntiAffinityGroup(t *testing.T) {
	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()
	clientFactory := plugin.NewClientFactory(tf.ExoscaleKey, tf.ExoscaleSecret)
	client := clientFactory.GetExoscaleClient()

	aas, err := client.RequestWithContext(context.Background(), &egoscale.ListAffinityGroups{})
	if err != nil {
		assert.Fail(t, "failed to list affinity groups (%v)", err)
	}
	assert.Equal(t, 1, len(aas.(*egoscale.ListAffinityGroupsResponse).AffinityGroup), "invalid number of AA's returned after initialization (%d)", len(aas.(*egoscale.ListAffinityGroupsResponse).AffinityGroup))

	i := aa.New()
	err = i.Run(clientFactory, context.Background())
	if err != nil {
		t.Fail()
	}

	aas, err = client.RequestWithContext(context.Background(), &egoscale.ListAffinityGroups{})
	if err != nil {
		assert.Fail(t, "failed to list affinity groups (%v)", err)
	}
	assert.Equal(t, 0, len(aas.(*egoscale.ListAffinityGroupsResponse).AffinityGroup), "invalid number of AA's returned after initialization (%d)", len(aas.(*egoscale.ListAffinityGroupsResponse).AffinityGroup))
}
