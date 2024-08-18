package user

import (
	"context"

	"github.com/npushpakumara/go-backend-template/internal/features/user/dto"
	"github.com/npushpakumara/go-backend-template/internal/features/user/entity"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

// UserService defines the methods that our User Service should implement.
type UserService interface {
	Register(ctx context.Context, user *dto.RegisterRequestDto) (*dto.UserResponseDto, error)
	UpdateUser(ctx context.Context, userId string, updates map[string]interface{}) error
}

// userServiceImpl is the concrete implementation of the UserService interface.
type userServiceImpl struct {
	userRepository UserRepository
}

// NewUserService creates a new instance of userServiceImpl with the provided UserRepository.
// This function initializes the user service with the repository it will use for data operations.
func NewUserService(userRepository UserRepository) UserService {
	return &userServiceImpl{userRepository}
}

// Register handles the registration of a new user.
// It takes a context and a RegisterRequestDto containing user details,
// hashes the user's password, and then inserts the user into the repository.
// If successful, it returns a UserResponseDto with the user's details; otherwise, it returns an error.
func (us *userServiceImpl) Register(ctx context.Context, user *dto.RegisterRequestDto) (*dto.UserResponseDto, error) {
	logger := logging.FromContext(ctx)
	requestBody := &entity.User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		logger.Errorw("user.service.Register failed to hash password: ", err)
		return nil, err
	}

	requestBody.Password = hashedPassword

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
