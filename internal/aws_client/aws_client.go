package awsclient

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/ses"
)

// Define a global variable to hold the AWSClient instance
// and a sync.Once variable to ensure the client is created only once.
var (
	client *AWSClient
	once   sync.Once
)

// AWSClient wraps the AWS Service's clients
type AWSClient struct {
	ses *ses.Client
}

// NewAWSClient initializes a new AWSClient instance with the specified AWS region.
// It uses sync.Once to ensure that the client is created only once, even if called concurrently.
func NewAWSClient(region string) *AWSClient {
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Fatalf("unable to load AWS SDK config, %v", err)
		}
		client = &AWSClient{
			ses: ses.NewFromConfig(cfg),
		}
	})

	return client
}

// GetSESClient returns the SES client from the AWSClient instance.
// This allows access to SES functionality for sending emails, etc.
func (c *AWSClient) GetSESClient() *ses.Client {
	return c.ses
}
