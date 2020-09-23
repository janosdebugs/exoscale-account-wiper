package sos_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"github.com/janoszen/exoscale-account-wiper/sos"
	"github.com/janoszen/exoscale-account-wiper/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemovingSOSBuckets(t *testing.T) {
	tf := terraform.New(t, "testdata")
	if tf == nil {
		// No Terraform integration available
		return
	}
	tf.Apply()
	defer tf.Destroy()
	clientFactory := plugin.NewClientFactory(tf.ExoscaleKey, tf.ExoscaleSecret)

	zoneName := "at-vie-1"
	client, err := clientFactory.GetS3Client(zoneName)
	if err != nil {
		assert.Fail(t, "failed to create S3 client (%v)", err)
	}
	listBucketsOutput, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		assert.Fail(t, "failed to list S3 buckets (%v)", err)
	}
	assert.Equal(t, 1, len(listBucketsOutput.Buckets), fmt.Sprintf("invalid number of SOS buckets returned (%d)", len(listBucketsOutput.Buckets)))

	securityGroup := sos.New()
	err = securityGroup.Run(clientFactory, context.Background())
	if err != nil {
		t.Fail()
	}

	listBucketsOutput, err = client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		assert.Fail(t, "failed to list S3 buckets (%v)", err)
	}
	assert.Equal(t, 0, len(listBucketsOutput.Buckets), fmt.Sprintf("invalid number of SOS buckets returned (%d)", len(listBucketsOutput.Buckets)))

}
