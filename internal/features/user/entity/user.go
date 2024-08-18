package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system.
// The struct fields are annotated with GORM tags to specify database constraints.
type User struct {
	*gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName   string    `gorm:"size:100;not null"`
	LastName    string    `gorm:"size:100"`
	Email       string    `gorm:"size:100;unique;not null"`
	Password    string    `gorm:"size:255"`
	PhoneNumber string    `gorm:"size:20"`
	IsActive    bool      `gorm:"type:boolean"`
}

// BeforeCreate is a GORM hook that is triggered before a new record is created in the database.
// It sets the ID field to a new UUID if it hasn't been set already.
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return
}
