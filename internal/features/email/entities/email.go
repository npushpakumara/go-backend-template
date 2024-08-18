package entities

type Email struct {
	From     string
	To       []string
	Subject  string
	Template string
	Data     map[string]string
}

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
