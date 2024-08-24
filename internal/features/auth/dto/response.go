package dto

// SignUpResponseDto is a Data Transfer Object (DTO) used to structure the response for a sign-up or any related action.
// It includes a status and a message, which provide feedback about the outcome of the operation.
type SignUpResponseDto struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// OAuthResponseDto is a Data Transfer Object (DTO) used to represent the user data returned after successful OAuth authentication.
// It includes essential user information such as ID, name, email, and OAuth provider details.
type OAuthResponseDto struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
}
