package plugin

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/exoscale/egoscale"
	apiv2 "github.com/exoscale/egoscale/api/v2"
)

func NewClientFactory(apiKey string, apiSecret string) *ClientFactory {
	return &ClientFactory{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

type ClientFactory struct {
	apiKey    string
	apiSecret string
}

func (c *ClientFactory) GetExoscaleClient() *egoscale.Client {
	return egoscale.NewClient("https://api.exoscale.ch/v1", c.apiKey, c.apiSecret)
}

func (c *ClientFactory) GetDnsClient() *egoscale.Client {
	return egoscale.NewClient("https://api.exoscale.ch/dns", c.apiKey, c.apiSecret)
}

func (c *ClientFactory) GetExoscaleV2Context(zoneName string, ctx context.Context) context.Context {
	return apiv2.WithEndpoint(ctx, apiv2.NewReqEndpoint("", zoneName))
}

func (c *ClientFactory) GetS3Client(zoneName string) (*s3.S3, error) {
	endpoint := fmt.Sprintf("https://sos-%s.exo.io", zoneName)
	s3Session, err := session.NewSession(&aws.Config{
		Region:   aws.String(zoneName),
		Endpoint: &endpoint,
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     c.apiKey,
				SecretAccessKey: c.apiSecret,
			},
		}),
	})
	if err != nil {
		return nil, err
	}
	return s3.New(s3Session), nil
}
