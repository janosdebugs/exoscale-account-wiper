package sshkeys_test

import (
	"context"
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"github.com/janoszen/exoscale-account-wiper/sshkeys"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingSshKeys(t *testing.T) {
	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()
	clientFactory := plugin.NewClientFactory(tf.ExoscaleKey, tf.ExoscaleSecret)

	v1Client := clientFactory.GetExoscaleClient()

	sshKeys, err := v1Client.ListWithContext(context.Background(), &egoscale.SSHKeyPair{})
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed to list SSH keys (%v)", err))
	}
	assert.Equal(t, 1, len(sshKeys), fmt.Sprintf("invalid number of SSH Keys returned (%d)", len(sshKeys)))

	sshKeyDeleter := sshkeys.New()
	err = sshKeyDeleter.Run(clientFactory, context.Background())
	if err != nil {
		t.Fail()
	}

	sshKeys, err = v1Client.ListWithContext(context.Background(), &egoscale.SSHKeyPair{})
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed to list SSH keys (%v)", err))
	}
	assert.Equal(t, 0, len(sshKeys), fmt.Sprintf("invalid number of SSH keys returned (%d)", len(sshKeys)))
}
