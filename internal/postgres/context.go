package postgres

import (
	"context"
	"gorm.io/gorm"
)

// contextKey is a custom type used to store and retrieve values in the context.
// Using a custom type helps to avoid any conflicts with other keys in the context.
type contextKey string

// dbKey is the key we use to store and retrieve the *gorm.DB instance (our database connection) in the context.
// It's of type contextKey, ensuring it's unique.
var dbKey = contextKey("db")

// WithDB adds a *gorm.DB instance (database connection) to the given context.
// This allows us to pass the context around in our application, and wherever we have the context, we can access the database connection.
func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

// FromContext retrieves the *gorm.DB instance (database connection) from the context.
// If the database connection is not found in the context, it returns the provided defaultDB.
func FromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	// If the context is nil, return the default database connection.
	if ctx == nil {
		return defaultDB
	}

	// Try to get the database connection from the context.
	if db, ok := ctx.Value(dbKey).(*gorm.DB); ok {
		// If found, return it.
		return db
	}

	// If not found, return the default database connection.
	return defaultDB
}
