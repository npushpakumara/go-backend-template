package auth

import (
	"context"
	"fmt"

	"github.com/npushpakumara/go-backend-template/internal/config"
	"github.com/npushpakumara/go-backend-template/internal/features/auth/dto"
	"github.com/npushpakumara/go-backend-template/internal/features/auth/tokens"
	"github.com/npushpakumara/go-backend-template/internal/features/email"
	"github.com/npushpakumara/go-backend-template/internal/features/email/entities"
	"github.com/npushpakumara/go-backend-template/internal/features/user"
	userDto "github.com/npushpakumara/go-backend-template/internal/features/user/dto"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

// AuthService defines the methods that our authentication service will implement.
type AuthService interface {
	// RegisterUser handles the user registration process, including saving user data
	// and sending a verification email.
	RegisterUser(ctx context.Context, user *dto.SignUpRequestDto) error

	// ActivateAccount activates a user account using the provided token.
	// This typically involves verifying the token and updating the user's status.
	ActivateAccount(ctx context.Context, token string) error
}

// authServiceImpl is a concrete implementation of the AuthService interface.
type authServiceImpl struct {
	userService  user.UserService   // Service responsible for user operations
	emailService email.EmailService // Service responsible for sending emails
	cfg          *config.Config     // Configuration settings for the application
}

// NewAuthService creates a new instance of authServiceImpl with the provided services and configuration.
// This function returns an AuthService interface that uses the authServiceImpl implementation.
func NewAuthService(userService user.UserService, emailService email.EmailService, cfg *config.Config) AuthService {
	return &authServiceImpl{userService, emailService, cfg}
}

// RegisterUser processes the registration of a new user. It converts the provided sign-up request
// data into a format suitable for the user service, registers the user, and sends a verification email.
// Returns an error if any step of the process fails.
func (as *authServiceImpl) RegisterUser(ctx context.Context, requestBody *dto.SignUpRequestDto) error {
	// Convert the sign-up request data to the format needed by the user service.
	userPayload := &userDto.RegisterRequestDto{
		FirstName:   requestBody.FirstName,
		LastName:    requestBody.LastName,
		Email:       requestBody.Email,
		Password:    requestBody.Password,
		PhoneNumber: requestBody.PhoneNumber,
	}

	// Register the user with the user service.
	newUser, err := as.userService.Register(ctx, userPayload)
	if err != nil {
		return err
	}

	// Send an account verification email to the newly registered user.
	if err := as.sendAccountVerificationEmail(ctx, newUser); err != nil {
		return err
	}

	return nil
}

// ActivateAccount activates a user account using the provided token.
// The token is used to find and update the user's status to active.
// Returns an error if token extraction or user update fails.
func (as *authServiceImpl) ActivateAccount(ctx context.Context, token string) error {
	logger := logging.FromContext(ctx)

	// Extract the user ID from the token.
	id, err := tokens.ExtractSubjectFromToken(as.cfg.JWT.Secret, token)
	if err != nil {
		logger.Errorw("auth.service.ActivateAccount failed to extract id from token")
		return err
	}

	// Prepare the payload to update the user's status.
	payload := map[string]interface{}{
		"is_active": true,
	}

	err = as.userService.UpdateUser(ctx, id, payload)
	if err != nil {
		return err
	}

	return nil
}

// sendAccountVerificationEmail creates a JWT token for account verification and sends an email to the user.
// The email contains a verification link with the token.
// Returns an error if token creation or email sending fails.
func (as *authServiceImpl) sendAccountVerificationEmail(ctx context.Context, requestBody *userDto.UserResponseDto) error {
	logger := logging.FromContext(ctx)

	// Create a new JWT token for account verification.
	tokenString, err := tokens.NewJwtToken(requestBody.ID, as.cfg.JWT.Secret, 172800)
	if err != nil {
		logger.Errorw("auth.service.sendAccountVerificationEmail failed to create jwt token: %v", err)
		return err // Return error if token creation fails.
	}

	// Prepare the email content.
	newEmail := &entities.Email{
		To:       []string{requestBody.Email},
		From:     as.cfg.AWS.SESConfig.EmailFrom,
		Subject:  entities.EmailTemplates["UserVerification"].Subject,
		Template: entities.EmailTemplates["UserVerification"].Template,
		Data: map[string]string{
			"name": requestBody.FirstName,
			"url":  fmt.Sprintf("%s/api/v1/auth/verify?token=%s", as.cfg.Server.Domain, tokenString),
		},
	}

	// Send the verification email using the email service.
	if err := as.emailService.SendEmail(ctx, *newEmail); err != nil {
		return err
	}

	return nil
}
