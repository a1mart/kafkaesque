package validators

import (
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Apply validation rules based on tags
func applyRule(rule string, field reflect.Value, fieldName string) error {
	switch {
	case strings.HasPrefix(rule, "min="):
		return validateMin(rule, field, fieldName)
	case strings.HasPrefix(rule, "max="):
		return validateMax(rule, field, fieldName)
	case strings.HasPrefix(rule, "length="):
		return validateLength(rule, field, fieldName)
	case rule == "required":
		return validateRequired(field, fieldName)
	case rule == "email":
		return validateEmail(field, fieldName)
	case strings.HasPrefix(rule, "regex="):
		return validateRegex(rule, field, fieldName)
	case strings.HasPrefix(rule, "enum="):
		return validateEnum(rule, field, fieldName)
	case strings.HasPrefix(rule, "range="):
		bounds := strings.TrimPrefix(rule, "range=")
		parts := strings.Split(bounds, ":")
		if len(parts) == 2 {
			min, _ := strconv.Atoi(parts[0])
			max, _ := strconv.Atoi(parts[1])
			value, _ := strconv.Atoi(field.String())
			if value < min || value > max {
				return fmt.Errorf("%s must be between %d and %d", fieldName, min, max)
			}
		}
	case rule == "url":
		return validateURL(field, fieldName)
	case rule == "uri":
		return validateURI(field, fieldName)
	case rule == "ip":
		return validateIPAddress(field, fieldName)
	case rule == "filepath":
		return validateFilePath(field, fieldName)
	}
	return nil
}

// Validation functions

func validateMin(rule string, field reflect.Value, fieldName string) error {
	min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
	if len(field.String()) < min {
		return fmt.Errorf("%s must be at least %d characters long", fieldName, min)
	}
	return nil
}

func validateMax(rule string, field reflect.Value, fieldName string) error {
	max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
	if len(field.String()) > max {
		return fmt.Errorf("%s must be at most %d characters long", fieldName, max)
	}
	return nil
}

// Example: validate:"length=8" (fixed length)
func validateLength(rule string, field reflect.Value, fieldName string) error {
	if strings.HasPrefix(rule, "length=") {
		length, _ := strconv.Atoi(strings.TrimPrefix(rule, "length="))
		if len(field.String()) != length {
			return fmt.Errorf("%s must be exactly %d characters long", fieldName, length)
		}
	}
	return nil
}

func validateRequired(field reflect.Value, fieldName string) error {
	if field.String() == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func validateEmail(field reflect.Value, fieldName string) error {
	_, err := mail.ParseAddress(field.String())
	if err != nil {
		return fmt.Errorf("%s must be a valid email", fieldName)
	}
	return nil
}

func validateRegex(rule string, field reflect.Value, fieldName string) error {
	parts := strings.Split(rule, ",")
	pattern := strings.TrimPrefix(parts[0], "regex=")
	options := ""
	if len(parts) > 1 && strings.HasPrefix(parts[1], "regexOptions=") {
		options = strings.TrimPrefix(parts[1], "regexOptions=")
	}

	regexFlags := ""
	if strings.Contains(options, "i") { //case insensitive
		regexFlags += "(?i)"
	}
	if strings.Contains(options, "m") { //multiline matching
		regexFlags += "(?m)"
	}

	fullPattern := regexFlags + pattern
	matched, _ := regexp.MatchString(fullPattern, field.String())
	if !matched {
		return fmt.Errorf("%s must match the pattern: %s", fieldName, pattern)
	}
	return nil
}

// Example: validate:"enum=red|green|blue"
func validateEnum(rule string, field reflect.Value, fieldName string) error {
	allowedValues := strings.Split(strings.TrimPrefix(rule, "enum="), "|")
	for _, value := range allowedValues {
		if field.String() == value {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of %s", fieldName, strings.Join(allowedValues, ", "))
}

// Example: validate:"url"
func validateURL(field reflect.Value, fieldName string) error {
	_, err := url.ParseRequestURI(field.String())
	if err != nil {
		return fmt.Errorf("%s must be a valid URL", fieldName)
	}
	return nil
}

func validateURI(field reflect.Value, fieldName string) error {
	parsedURI, err := url.Parse(field.String())
	if err != nil || parsedURI.Scheme == "" || parsedURI.Host == "" {
		return fmt.Errorf("%s must be a valid URI", fieldName)
	}
	return nil
}

// Example: validate:"ip"
func validateIPAddress(field reflect.Value, fieldName string) error {
	if net.ParseIP(field.String()) == nil {
		return fmt.Errorf("%s must be a valid IP address", fieldName)
	}
	return nil
}

// Example: validate:"filepath"
func validateFilePath(field reflect.Value, fieldName string) error {
	if _, err := os.Stat(field.String()); os.IsNotExist(err) {
		return fmt.Errorf("%s is not a valid file path", fieldName)
	}
	return nil
}
