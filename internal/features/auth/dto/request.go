package dto

// SignUpRequestDto is a Data Transfer Object (DTO) used to capture and validate the data required for a new user sign-up.
// It includes fields for the user's first and last names, email, password, and phone number, all of which are required.
type SignUpRequestDto struct {
	FirstName   string `json:"first_name" binding:"required,min=2,max=100"`
	LastName    string `json:"last_name" binding:"required,min=2,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
	PhoneNumber string `json:"phone_number" binding:"required,e164,min=12,max=12"`
}

// SignInRequestDto is a Data Transfer Object (DTO) used to capture and validate the data required for user sign-in.
// It includes the user's email and password, both of which are required.
type SignInRequestDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

// PasswordResetRequestDto is a Data Transfer Object (DTO) used to capture and validate the data required for a password reset.
// It includes the user's email, current password, and new password, all of which are required.
type PasswordResetRequestDto struct {
	Email           string `json:"email" binding:"required,email"`
	CurrentPassword string `json:"current_password" binding:"required,min=8,max=100"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=100"`
}
