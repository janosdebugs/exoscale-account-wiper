package sg_test

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/sg"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingSecurityGroups(t *testing.T) {
	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()

	v1Client := egoscale.NewClient("https://api.exoscale.ch/v1", tf.ExoscaleKey, tf.ExoscaleSecret)
	sgPrototype := &egoscale.SecurityGroup{}
	sgs, err := v1Client.ListWithContext(context.Background(), sgPrototype)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed to list security groups (%v)", err))
	}
	assert.Equal(t, 2, len(sgs), fmt.Sprintf("invalid number of security groups returned (%d)", len(sgs)))

	securityGroup := sg.New()
	err = securityGroup.Run(v1Client, context.Background())
	if err != nil {
		t.Fail()
	}

	sgs, err = v1Client.ListWithContext(context.Background(), sgPrototype)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed to list security groups (%v)", err))
	}
	assert.Equal(t, 1, len(sgs), fmt.Sprintf("invalid number of security groups returned (%d)", len(sgs)))
}
