package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/npushpakumara/go-backend-template/internal/config"
	"gorm.io/driver/postgres"

	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

// NewDatabase creates and configures a new database connection using GORM.
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Initialize variables to hold the database connection, error, and logger
	var (
		db  *gorm.DB
		err error
		// Create a custom logger for GORM using the zap
		logger = NewLogger(time.Second, true, zapcore.Level(cfg.DB.LogLevel))
	)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)

	// Attempt to connect to the database up to 10 times with retries
	for i := 0; i < 10; i++ {
		// Try to open a database connection with GORM using the Postgres driver
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger, // Use the custom logger for GORM logging
		})

		// If the connection was successful, exit the loop
		if err == nil {
			break
		}

		// Log the error and retry after a short delay
		log.Printf("Attempt %d: Failed to connect to the database: %v", i+1, err)
		time.Sleep(500 * time.Millisecond)
	}

	// If we failed to connect after all attempts, return the error
	if err != nil {
		return nil, err
	}

	// Get the underlying SQL database connection from GORM
	pgDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configure the connection pool settings
	pgDB.SetMaxOpenConns(cfg.DB.Pool.MaxOpen)        // Maximum number of open connections to the database
	pgDB.SetMaxIdleConns(cfg.DB.Pool.MaxIdle)        // Maximum number of idle connections in the pool
	pgDB.SetConnMaxLifetime(cfg.DB.Pool.MaxLifetime) // Maximum lifetime of a connection before it is reused

	err = db.Exec("CREATE SCHEMA IF NOT EXISTS auc").Error
	if err != nil {
		return nil, err
	}

	// If database migration is enabled, run migrations
	if cfg.DB.Migrations {
		// Call the migrateDB function to apply migrations
		err := migrateAndSeed(db)
		if err != nil {
			return nil, err
		}
	}

	// Return the successfully connected and configured GORM database instance
	return db, nil
}
