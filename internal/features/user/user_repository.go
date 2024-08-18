package user

import (
	"context"
	"errors"

	"github.com/npushpakumara/go-backend-template/internal/features/user/entity"
	"github.com/npushpakumara/go-backend-template/internal/postgres"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserRepository defines the interface for user-related data operations.
type UserRepository interface {
	// Insert adds a new user to the database.
	// It returns the inserted user and an error if something goes wrong.
	Insert(ctx context.Context, user *entity.User) (*entity.User, error)

	// FindByEmail retrieves a user by their email address.
	// It returns the user if found or an error if something goes wrong or the user does not exist.
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByID retrieves a user by their unique identifier (ID).
	// It returns the user if found or an error if something goes wrong or the user does not exist.
	FindByID(ctx context.Context, id string) (*entity.User, error)

	// Update modifies the details of an existing user identified by ID.
	// It takes a map of field names and values to update and returns an error if the update fails.
	Update(ctx context.Context, id string, updates map[string]interface{}) error
}

// userRepositoryImpl is a concrete implementation of the UserRepository interface.
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of userRepositoryImpl with the provided database connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db}
}

// Insert adds a new user to the database.
// It logs the operation and handles potential errors, including checking for duplicate entries.
func (us *userRepositoryImpl) Insert(ctx context.Context, user *entity.User) (*entity.User, error) {
	logger := logging.FromContext(ctx)
	db := postgres.FromContext(ctx, us.db)

	logger.Debugw("user.db.Insert", "user", user)
	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		if pgErr := postgres.IsPgxError(err); errors.Is(pgErr, postgres.ErrKeyDuplicate) {
			logger.Warn("user.db.Insert user already exists")
			return nil, postgres.ErrKeyDuplicate
		}
		logger.Errorw("user.db.Insert failed to save: %v", err)
		return nil, err
	}
	return user, nil
}

// FindByEmail searches for a user based on their email address.
// It logs the search operation and handles errors, including the case where the user is not found.
func (us *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	logger := logging.FromContext(ctx)
	db := postgres.FromContext(ctx, us.db)

	logger.Debugw("user.db.FindByEmail", "email", email)

	var user *entity.User
	if err := db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("user.db.FindByEmail user not found")
			return nil, postgres.ErrRecordNotFound
		}
		logger.Errorw("user.db.FindByEmail failed to find user: %v", err)
		return nil, err
	}
	return user, nil
}

// FindByID retrieves a user based on their ID.
// It logs the search operation and handles errors, including the case where the user is not found.
func (us *userRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.User, error) {
	logger := logging.FromContext(ctx)
	db := postgres.FromContext(ctx, us.db)

	logger.Debugw("user.db.FindByID", "id", id)

	var user *entity.User
	if err := db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("user.db.FindByID user not found")
			return nil, postgres.ErrRecordNotFound
		}
		logger.Errorw("user.db.FindByID failed to find user: %v", err)
		return nil, err
	}
	return user, nil
}

// Update modifies an existing user's details based on their ID.
// It logs the update operation and handles errors, including the case where the user is not found.
func (us *userRepositoryImpl) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	logger := logging.FromContext(ctx)
	db := postgres.FromContext(ctx, us.db)

	logger.Debugw("user.db.Update", id, updates)

	var user *entity.User
	if err := db.WithContext(ctx).Model(&user).Clauses(clause.Returning{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("user.db.Update user not found")
			return postgres.ErrRecordNotFound
		}
		logger.Errorw("user.db.Update failed to update user: %v", err)
		return err
	}

	return nil
}
