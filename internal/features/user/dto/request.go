package dto

// RegisterRequestDto is a data transfer object used for handling
// user registration information. It contains fields that capture
// the essential details needed to register a new user.
type RegisterRequestDto struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber string
	Provider    string
	ProviderID  string
}
