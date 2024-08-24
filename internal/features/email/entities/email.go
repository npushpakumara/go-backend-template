package entities

// Email represents the structure of an email message.
type Email struct {
	From     string
	To       []string
	Subject  string
	Template string
	Data     map[string]string
}

// EmailTemplates is a map that stores predefined email templates with their subjects and template names.
// Each template is identified by a unique key, such as "UserVerification" or "PasswordReset".
var EmailTemplates = map[string]struct {
	Subject  string
	Template string
}{
	"UserVerification": {
		Subject:  "User Activation Email",
		Template: "email-verification",
	},
	"PasswordReset": {
		Subject:  "Password Reset Request",
		Template: "password-reset",
	},
}
