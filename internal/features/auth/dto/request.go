package dto

// SignUpRequestDto represents the data transfer object used for handling user sign-up requests.
// It contains fields for user information and their validation rules.
type SignUpRequestDto struct {
	FirstName   string `json:"first_name" binding:"required,min=2,max=100"`
	LastName    string `json:"last_name" binding:"required,min=2,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
	PhoneNumber string `json:"phone_number" binding:"required,e164,min=12,max=12"`
}
