package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/npushpakumara/go-backend-template/internal/config"
	"github.com/npushpakumara/go-backend-template/internal/features/email/entities"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

// smtpServiceImpl is an implementation of an email service that uses SMTP to send emails.
// It stores the SMTP server address and authentication details.
type smtpServiceImpl struct {
	Server string
	Auth   smtp.Auth
}

// NewSMTPEmailService initializes and returns a new instance of smtpServiceImpl.
// It sets up the SMTP authentication using the provided configuration and constructs the server address.
func NewSMTPEmailService(cfg *config.Config) Service {
	auth := smtp.PlainAuth("", cfg.Mail.SMTP.Username, cfg.Mail.SMTP.Password, cfg.Mail.SMTP.Server)
	return &smtpServiceImpl{
		Server: fmt.Sprintf("%s:%d", cfg.Mail.SMTP.Server, cfg.Mail.SMTP.Port),
		Auth:   auth,
	}
}

// SendEmail sends an email using the SMTP server specified in smtpServiceImpl.
// It logs any errors encountered during the sending process.
func (s *smtpServiceImpl) SendEmail(ctx context.Context, email entities.Email) error {
	logger := logging.FromContext(ctx)

	subject := "Subject: " + email.Subject + "\n"
	contentType := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + contentType + email.Data)

	err := smtp.SendMail(s.Server, s.Auth, email.From, email.To, msg)
	if err != nil {
		logger.Errorf("email.service.SendEmail error while sending email via Gmail: %w", err)
		return err
	}

	return nil
}
