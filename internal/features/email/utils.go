package email

import (
	"bytes"
	"fmt"
	"text/template"
)

// ParseTemplate parses a template string and applies the provided data to it, returning the resulting string.
// If there is an error during the parsing or execution of the template, it returns an empty string and the error.
func ParseTemplate(templateString string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("internal/features/email/templates/%s", templateString))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
