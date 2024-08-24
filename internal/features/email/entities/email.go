package entities

// Email represents the structure of an email message.
type Email struct {
	From    string
	To      []string
	Subject string
	Data    string
}

// VerificationEmailData is a struct that holds the dynamic data needed to populate a verification email template.
// It includes the recipient's name and a verification link, which will be inserted into the email template.
type VerificationEmailData struct {
	Name string
	Link string
}

// EmailTemplates is a map that stores predefined email templates with their subjects and template names.
// Each template is identified by a unique key, such as "UserVerification" or "PasswordReset".
var EmailTemplates = map[string]struct {
	Subject  string
	Template string
}{
	"UserVerification": {
		Subject:  "User Activation Email",
		Template: "account-verification.html",
	},
	"PasswordReset": {
		Subject:  "Password Reset Request",
		Template: "password-reset.html",
	},
}