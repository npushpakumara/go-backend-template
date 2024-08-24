package email

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	awsclient "github.com/npushpakumara/go-backend-template/internal/aws_client"
	"github.com/npushpakumara/go-backend-template/internal/features/email/entities"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

// Service defines an interface for sending emails.
// It provides a method to send an email with a given context and email details.
type Service interface {
	SendEmail(c context.Context, email entities.Email) error
}

// emailServiceImpl is a concrete implementation of the Service interface.
// It uses an AWS client to send emails through AWS SES (Simple Email Service).
type emailServiceImpl struct {
	AWSClient *awsclient.AWSClient
}

// NewSESEmailService creates a new instance of emailServiceImpl.
// It initializes the service with the given AWS client.
// This function returns an Service interface that wraps the emailServiceImpl.
func NewSESEmailService(awsClient *awsclient.AWSClient) Service {
	return &emailServiceImpl{
		AWSClient: awsClient,
	}
}

// SendEmail sends an email using AWS SES with the provided context and email details.
// It marshals the email data into JSON format and constructs the input for the SES API.
// If there is an error in marshalling the data or sending the email, it logs the error
// and returns it. Otherwise, it returns nil indicating success.
func (s *emailServiceImpl) SendEmail(ctx context.Context, email entities.Email) error {
	logger := logging.FromContext(ctx)

	jsonData, err := json.Marshal(email.Data)
	if err != nil {
		logger.Errorw("email.service.SendEmail failed to marshal template data: %v", err)
		return err
	}

	input := &ses.SendTemplatedEmailInput{
		Destination: &types.Destination{
			ToAddresses: email.To,
		},
		Source:       aws.String(email.From),
		Template:     aws.String(email.Template),
		TemplateData: aws.String(string(jsonData)),
	}

	_, err = s.AWSClient.GetSESClient().SendTemplatedEmail(ctx, input)
	if err != nil {
		logger.Errorw("email.service.SendEmail error while sending email via aws ses: %w", err)
		return err
	}
	return nil
}
