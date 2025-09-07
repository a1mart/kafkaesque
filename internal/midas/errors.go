package validators

import (
	"fmt"
	"strings"
)

/*
var Messages = map[string]string{
	"required": "This field is required.",
	"email":    "Please provide a valid email address.",
	"password": "Password must contain...",
}

func getMessage(key string, fieldName string) string {
	return fmt.Sprintf(Messages[key], fieldName)
}
*/

type ValidationError struct {
	Field   string
	Message string
}

type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) AddError(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// Error implements the error interface for ValidationResult.
// It formats all validation errors as a single error message.
func (r *ValidationResult) Error() string {
	if !r.HasErrors() {
		return ""
	}

	var messages []string
	for _, e := range r.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}
	return strings.Join(messages, "; ")
}
