package user

import (
	"context"

	"github.com/npushpakumara/go-backend-template/internal/features/user/dto"
	"github.com/npushpakumara/go-backend-template/internal/features/user/entity"
	"github.com/npushpakumara/go-backend-template/internal/postgres"
)

// UserService defines the methods that our User Service should implement.
type UserService interface {
	CreateUser(ctx context.Context, user *dto.RegisterRequestDto) (*dto.UserResponseDto, error)
	UpdateUser(ctx context.Context, userId string, updates map[string]interface{}) error
	GetUserByID(ctx context.Context, userID string) (*dto.UserResponseDto, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponseDto, error)
}

// userServiceImpl is the concrete implementation of the UserService interface.
type userServiceImpl struct {
	userRepository UserRepository
}

// NewUserService creates a new instance of userServiceImpl with the provided UserRepository.
// This function initializes the user service with the repository it will use for data operations.
func NewUserService(userRepository UserRepository, transactionManager postgres.TransactionManager) UserService {
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
		PhoneNumber: user.PhoneNumber,
		Provider:    user.Provider,
		ProviderID:  user.ProviderID,
	}

	if user.ProviderID == "" {
		requestBody.Password = user.Password
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
func (us *userServiceImpl) UpdateUser(ctx context.Context, userId string, updates map[string]interface{}) error {

	err := us.userRepository.Update(ctx, userId, updates)
	if err != nil {
		return err
	}
	return nil
}

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
