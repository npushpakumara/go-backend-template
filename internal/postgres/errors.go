package postgres

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

// Predefined error variables for specific PostgreSQL error codes.
var (
	ErrKeyDuplicate        = errors.New("duplicate key found")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrUniqueViolation     = errors.New("unique key violation")
	ErrRecordNotFound      = errors.New("record not found")
)

// IsPgxError checks if the given error is a PostgreSQL error and returns a corresponding custom error.
func IsPgxError(err error) error {
	if err == nil {
		return nil
	}

	// Check if the error is a PostgreSQL error.
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505":
			return ErrKeyDuplicate
		case "23503":
			return ErrForeignKeyViolation
		case "23514":
			return ErrUniqueViolation
		default:
			return errors.New("database error: " + pgErr.Message)
		}
	}

	return err
}
