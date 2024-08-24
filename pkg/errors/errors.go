package errors

import "errors"

// ErrInvalidToken is returned when the provided JWT token is invalid.
// This could happen if the token is malformed, expired, or fails verification.
var ErrInvalidToken = errors.New("invalid jwt token")

// ErrAccountNotActive is returned when a user attempts to perform an action,
// but their account is not active. This typically indicates that the user needs
// to activate their account before they can proceed.
var ErrAccountNotActive = errors.New("user is not active")

// ErrIncorrectPassword is returned when a user provides an incorrect password
// during authentication. This prevents unauthorized access to the account.
var ErrIncorrectPassword = errors.New("incorrect password")

// ErrEmailLinkedToOauth is returned when a user attempts to sign up or log in
// using an email that is already associated with an OAuth account. This error
// informs the user that they should use their OAuth provider to log in instead.
var ErrEmailLinkedToOauth = errors.New("email associated with oauth account")

// ErrorResponse represents the structure of an error response.
// It includes a status, a message, and optionally additional error details.
type ErrorResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}
