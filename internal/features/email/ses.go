package email

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	awsclient "github.com/npushpakumara/go-backend-template/internal/aws_client"
	"github.com/npushpakumara/go-backend-template/internal/features/email/entities"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

// sesEmailServiceImpl is a concrete implementation of the Service interface.
// It uses an AWS client to send emails through AWS SES (Simple Email Service).
type sesEmailServiceImpl struct {
	AWSClient *awsclient.AWSClient
}

// NewSESEmailService creates a new instance of emailServiceImpl.
// It initializes the service with the given AWS client.
// This function returns an Service interface that wraps the emailServiceImpl.
func NewSESEmailService(awsClient *awsclient.AWSClient) Service {
	return &sesEmailServiceImpl{
		AWSClient: awsClient,
	}
}

// SendEmail sends an email using AWS SES with the provided context and email details.
// It marshals the email data into JSON format and constructs the input for the SES API.
// If there is an error in marshalling the data or sending the email, it logs the error
// and returns it. Otherwise, it returns nil indicating success.
func (s *sesEmailServiceImpl) SendEmail(ctx context.Context, email entities.Email) error {
	logger := logging.FromContext(ctx)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: email.To,
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(email.Data),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(email.Subject),
			},
		},
		Source: aws.String(email.From),
	}

	_, err := s.AWSClient.GetSESClient().SendEmail(ctx, input)
	if err != nil {
		logger.Errorw("email.service.SendEmail error while sending email via aws ses: %w", err)
		return err
	}
	return nil
}
