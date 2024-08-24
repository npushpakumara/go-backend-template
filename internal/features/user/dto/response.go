package dto

import "time"

// UserResponseDto represents the data structure for a user's response.
// It contains all the information that will be sent back to the client when querying user details.
type UserResponseDto struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber string
	IsActive    bool
	Provider    string
	ProviderID  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
