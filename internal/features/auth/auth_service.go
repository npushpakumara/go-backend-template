package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/markbates/goth"
	"github.com/npushpakumara/go-backend-template/internal/config"
	"github.com/npushpakumara/go-backend-template/internal/features/auth/dto"
	"github.com/npushpakumara/go-backend-template/internal/features/auth/tokens"
	"github.com/npushpakumara/go-backend-template/internal/features/email"
	"github.com/npushpakumara/go-backend-template/internal/features/email/entities"
	"github.com/npushpakumara/go-backend-template/internal/features/user"
	userDto "github.com/npushpakumara/go-backend-template/internal/features/user/dto"
	"github.com/npushpakumara/go-backend-template/internal/postgres"
	apiError "github.com/npushpakumara/go-backend-template/pkg/errors"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

// Service defines the methods that our authentication service will implement.
type Service interface {
	// RegisterUser handles the process of registering a new user.
	// It accepts a SignUpRequestDto containing the user's registration details and performs necessary actions such as
	// validating the input, storing the user's data, and sending a confirmation email.
	RegisterUser(ctx context.Context, user *dto.SignUpRequestDto) error

	// LoginUser handles the user login process.
	// It accepts a SignInRequestDto containing the user's email and password, validates the credentials,
	// and returns the user's ID if successful. If login fails, it returns an appropriate error.
	LoginUser(ctx context.Context, request *dto.SignInRequestDto) (string, error)

	// ResetPassword handles the process of resetting a user's password.
	// It accepts a PasswordResetRequestDto containing the user's current and new passwords, verifies the current password,
	// and updates the user's password in the database if validation is successful.
	ResetPassword(ctx context.Context, request *dto.PasswordResetRequestDto) error

	// ActivateAccount handles the activation of a user's account.
	// It accepts a token string, verifies its validity, and activates the account associated with the token.
	// It returns the user's ID if activation is successful.
	ActivateAccount(ctx context.Context, token string) (string, error)

	// GetUserByID retrieves a user's details based on their ID.
	// It returns a UserResponseDto containing the user's information, or an error if the user is not found.
	GetUserByID(ctx context.Context, id string) (*userDto.UserResponseDto, error)

	// SendAccountVerificationEmail sends an account verification email to the user.
	// It accepts a UserResponseDto containing the user's details, generates a verification token,
	// and sends the email. Returns an error if the email cannot be sent.
	SendAccountVerificationEmail(ctx context.Context, requestBody *userDto.UserResponseDto) error

	// HandleOAuthUser handles the authentication of a user via OAuth.
	// It accepts a Goth User object containing the OAuth user's details, processes the user (e.g., linking accounts, creating a new user),
	// and returns an OAuthResponseDto with the necessary information, or an error if the process fails.
	HandleOAuthUser(ctx context.Context, gothUser goth.User) (*dto.OAuthResponseDto, error)
}

// authServiceImpl is a concrete implementation of the Service interface.
type authServiceImpl struct {
	userService        user.Service  // Service responsible for user operations
	emailService       email.Service // Service responsible for sending emails
	transactionManager postgres.TransactionManager
	cfg                *config.Config // Configuration settings for the application
}

// NewAuthService creates a new instance of authServiceImpl with the provided services and configuration.
// This function returns an Service interface that uses the authServiceImpl implementation.
func NewAuthService(userService user.Service, emailService email.Service, transactionManager postgres.TransactionManager, cfg *config.Config) Service {
	return &authServiceImpl{userService, emailService, transactionManager, cfg}
}

// RegisterUser processes the registration of a new user. It converts the provided sign-up request
// data into a format suitable for the user service, registers the user, and sends a verification email.
// Returns an error if any step of the process fails.
func (as *authServiceImpl) RegisterUser(c context.Context, requestBody *dto.SignUpRequestDto) error {
	logger := logging.FromContext(c)

	ctx, err := as.transactionManager.Begin(c)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			as.transactionManager.Rollback(ctx)
		}
	}()

	// Convert the sign-up request data to the format needed by the user service.
	userPayload := &userDto.RegisterRequestDto{
		FirstName:   requestBody.FirstName,
		LastName:    requestBody.LastName,
		Email:       requestBody.Email,
		Password:    requestBody.Password,
		PhoneNumber: requestBody.PhoneNumber,
	}

	hashedPassword, err := hashPassword(requestBody.Password)
	if err != nil {
		logger.Errorw("auth.service.RegisterUser failed to hash password: ", err)
		return err
	}

	userPayload.Password = hashedPassword

	// Register the user with the user service.
	newUser, err := as.userService.CreateUser(ctx, userPayload)
	if err != nil {
		return err
	}

	// Send an account verification email to the newly registered user.
	if err := as.SendAccountVerificationEmail(ctx, newUser); err != nil {
		return err
	}

	as.transactionManager.Commit(ctx)

	return nil
}

// ActivateAccount activates a user account using the provided token.
// The token is used to find and update the user's status to active.
// Returns an error if token extraction or user update fails.
func (as *authServiceImpl) ActivateAccount(ctx context.Context, token string) (string, error) {
	logger := logging.FromContext(ctx)

	// Extract the user ID from the token.
	id, err := tokens.ExtractSubjectFromToken(as.cfg.JWT.Secret, token)
	if err != nil {
		logger.Errorw("auth.service.ActivateAccount failed to extract id from token", err)
		return "", err
	}

	// Prepare the payload to update the user's status.
	payload := map[string]interface{}{
		"is_active": true,
	}

	err = as.userService.UpdateUser(ctx, id, payload)
	if err != nil {
		return "", err
	}

	return id, nil
}

// SendAccountVerificationEmail creates a JWT token for account verification and sends an email to the user.
// The email contains a verification link with the token.
// Returns an error if token creation or email sending fails.
func (as *authServiceImpl) SendAccountVerificationEmail(ctx context.Context, requestBody *userDto.UserResponseDto) error {
	logger := logging.FromContext(ctx)

	// Create a new JWT token for account verification.
	tokenString, err := tokens.NewJwtToken(requestBody.ID, as.cfg.JWT.Secret, time.Hour*48)
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

// HandleOAuthUser handles the process of registering a user via an OAuth provider.
// It takes in the OAuth user information, creates a user registration payload,
// and attempts to register the user using the userService.
func (as *authServiceImpl) HandleOAuthUser(ctx context.Context, gothUser goth.User) (*dto.OAuthResponseDto, error) {
	userPayload := &userDto.RegisterRequestDto{
		FirstName:  gothUser.FirstName,
		LastName:   gothUser.LastName,
		Email:      gothUser.Email,
		Provider:   gothUser.Provider,
		ProviderID: gothUser.UserID,
	}

	resp, err := as.userService.CreateUser(ctx, userPayload)
	if err != nil {
		if errors.Is(err, postgres.ErrKeyDuplicate) {
			resp, err = as.userService.GetUserByEmail(ctx, gothUser.Email)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &dto.OAuthResponseDto{
		ID:         resp.ID,
		FirstName:  resp.FirstName,
		LastName:   resp.LastName,
		Email:      resp.Email,
		Provider:   resp.Provider,
		ProviderID: resp.ProviderID,
	}, nil
}

// GetUserByID retrieves a user by their ID and returns a UserResponseDto.
// It logs any errors that occur during the process.
func (as *authServiceImpl) GetUserByID(ctx context.Context, id string) (*userDto.UserResponseDto, error) {
	logger := logging.FromContext(ctx)

	user, err := as.userService.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("auth.service.GetUserByID failed to get user by id: %v", err)
		return nil, err
	}

	return user, nil
}

// LoginUser attempts to log in a user based on the provided SignInRequestDto.
// It performs various checks such as validating the email, checking if the account is active, and verifying the password.
func (as *authServiceImpl) LoginUser(ctx context.Context, requestBody *dto.SignInRequestDto) (string, error) {
	logger := logging.FromContext(ctx)

	resp, err := as.userService.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		logger.Errorf("auth.service.LoginUser failed to get user by email: %v", err)
		return "", err
	}

	if resp.ProviderID != "" {
		logger.Errorw("auth.service.LoginUser failed to login", "email associate with oauth account")
		return "", apiError.ErrEmailLinkedToOauth
	}

	if !resp.IsActive {
		logger.Errorf("auth.service.LoginUser account is not activated")
		return "", apiError.ErrAccountNotActive
	}

	if err := checkPassword(resp.Password, requestBody.Password); err != nil {
		if errors.Is(err, apiError.ErrIncorrectPassword) {
			logger.Errorw("auth.service.LoginUser failed to login", "invalid password", err)
			return "", err
		}
		return "", err
	}

	return resp.ID, nil
}

// ResetPassword allows a user to reset their password by providing the current and new passwords.
// It first verifies the current password and then updates the user's password in the database.
func (as *authServiceImpl) ResetPassword(ctx context.Context, request *dto.PasswordResetRequestDto) error {
	logger := logging.FromContext(ctx)

	resp, err := as.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		logger.Errorf("auth.service.ResetPassword failed to get user by email: %v", err)
		return err
	}

	err = checkPassword(resp.Password, request.CurrentPassword)
	if err != nil {
		logger.Errorf("auth.service.ResetPassword incorrect current password: %v", err)
		return apiError.ErrIncorrectPassword
	}

	hashedPassword, err := hashPassword(request.NewPassword)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"password": hashedPassword,
	}

	err = as.userService.UpdateUser(ctx, resp.ID, payload)
	if err != nil {
		return err
	}

	return nil
}
