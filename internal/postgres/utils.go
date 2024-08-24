package postgres

import (
	"log"

	"github.com/npushpakumara/go-backend-template/internal/features/user/entity"
	"gorm.io/gorm"
)

// migrateAndSeed is a function that performs database migration and seeding.
func migrateAndSeed(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
		return err
	}
	return nil
}
