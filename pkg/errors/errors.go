package errors

import "errors"

var ErrInvalidToken = errors.New("invalid jwt token")

// ErrorResponse represents the structure of an error response.
// It includes a status, a message, and optionally additional error details.
type ErrorResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}
