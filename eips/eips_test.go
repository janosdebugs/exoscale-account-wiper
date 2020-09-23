package eips_test

import (
	"context"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/eips"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingElasticIp(t *testing.T) {
	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()
	clientFactory := plugin.NewClientFactory(tf.ExoscaleKey, tf.ExoscaleSecret)
	client := clientFactory.GetExoscaleClient()

	req := egoscale.IPAddress{}
	ips, err := client.ListWithContext(context.Background(), &req)
	if err != nil {
		assert.Fail(t, "failed list EIP's after initialization (%v)", err)
		return
	}
	assert.Equal(t, 1, len(ips), "invalid number of EIP's returned after initialization (%d)", len(ips))

	i := eips.New()
	err = i.Run(clientFactory, context.Background())
	if err != nil {
		t.Fail()
	}

	req = egoscale.IPAddress{}
	ips, err = client.ListWithContext(context.Background(), &req)
	if err != nil {
		assert.Fail(t, "failed list EIP's after run (%v)", err)
		return
	}
	assert.Equal(t, 0, len(ips), "invalid number of EIP's returned after run (%d)", len(ips))

}
