package email

import (
	"context"

	awsclient "github.com/npushpakumara/go-backend-template/internal/aws_client"
	"github.com/npushpakumara/go-backend-template/internal/config"
	"github.com/npushpakumara/go-backend-template/internal/features/email/entities"
)

// Service defines an interface for sending emails.
// It provides a method to send an email with a given context and email details.
type Service interface {
	SendEmail(c context.Context, email entities.Email) error
}

// Provider defines the available email providers.
type Provider string

const (
	providerSES  Provider = "ses"
	providerSMTP Provider = "smtp"
)

// NewEmailService creates a new email service based on the given provider.
func NewEmailService(cfg *config.Config, awsClient *awsclient.AWSClient) Service {
	switch Provider(cfg.Mail.Provider) {
	case providerSES:
		return NewSESEmailService(awsClient)
	case providerSMTP:
		return NewSMTPEmailService(cfg)
	default:
		return nil
	}
}
