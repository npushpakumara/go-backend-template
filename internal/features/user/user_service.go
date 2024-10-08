package user

import (
	"context"

	"github.com/npushpakumara/go-backend-template/internal/features/user/dto"
	"github.com/npushpakumara/go-backend-template/internal/features/user/entity"
	"github.com/npushpakumara/go-backend-template/internal/postgres"
)

// Service defines the methods that our User Service should implement.
type Service interface {
	CreateUser(ctx context.Context, user *dto.RegisterRequestDto) (*dto.UserResponseDto, error)
	UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error
	GetUserByID(ctx context.Context, userID string) (*dto.UserResponseDto, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponseDto, error)
}

// userServiceImpl is the concrete implementation of the Service interface.
type userServiceImpl struct {
	userRepository Repository
}

// NewUserService creates a new instance of userServiceImpl with the provided Repository.
// This function initializes the user service with the repository it will use for data operations.
func NewUserService(userRepository Repository, transactionManager postgres.TransactionManager) Service {
	return &userServiceImpl{userRepository}
}

// CreateUser handles the registration of a new user.
// It takes a context and a RegisterRequestDto containing user details,
// hashes the user's password, and then inserts the user into the repository.
// If successful, it returns a UserResponseDto with the user's details; otherwise, it returns an error.
func (us *userServiceImpl) CreateUser(ctx context.Context, user *dto.RegisterRequestDto) (*dto.UserResponseDto, error) {

	requestBody := &entity.User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		Provider:    user.Provider,
		ProviderID:  user.ProviderID,
	}

	// If the user is not an oauth user, then set the password
	if user.ProviderID != "" {
		requestBody.Password = ""
		requestBody.IsActive = true
	}

	newUser, err := us.userRepository.Insert(ctx, requestBody)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponseDto{
		ID:        newUser.ID.String(),
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

// UpdateUser updates the details of an existing user based on the userId and the updates map.
func (us *userServiceImpl) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {

	err := us.userRepository.Update(ctx, userID, updates)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user by their ID and returns a UserResponseDto containing the user's details.
// It first fetches the user from the repository using the user ID, then maps the user entity to a UserResponseDto.
func (us *userServiceImpl) GetUserByID(ctx context.Context, userID string) (*dto.UserResponseDto, error) {
	user, err := us.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userDto := &dto.UserResponseDto{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		IsActive:  user.IsActive,
	}
	return userDto, nil
}

// GetUserByEmail retrieves a user by their email and returns a UserResponseDto containing the user's details.
// It first fetches the user from the repository using the email, then maps the user entity to a UserResponseDto.
func (us *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponseDto, error) {
	user, err := us.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	userDto := &dto.UserResponseDto{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Password:   user.Password,
		CreatedAt:  user.CreatedAt,
		IsActive:   user.IsActive,
		ProviderID: user.ProviderID,
	}
	return userDto, nil
}
