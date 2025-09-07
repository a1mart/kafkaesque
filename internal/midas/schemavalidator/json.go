package schemavalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// JSONSchema represents a JSON Schema definition
type JSONSchema struct {
	Type        interface{}            `json:"type"`
	Properties  map[string]*JSONSchema `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Enum        []interface{}          `json:"enum,omitempty"`
	Minimum     *float64               `json:"minimum,omitempty"`
	Maximum     *float64               `json:"maximum,omitempty"`
	MinLength   *int                   `json:"minLength,omitempty"`
	MaxLength   *int                   `json:"maxLength,omitempty"`
	Pattern     *string                `json:"pattern,omitempty"`
	Format      *string                `json:"format,omitempty"`
	Items       *JSONSchema            `json:"items,omitempty"`
	MinItems    *int                   `json:"minItems,omitempty"`
	MaxItems    *int                   `json:"maxItems,omitempty"`
	UniqueItems bool                   `json:"uniqueItems,omitempty"`
}

// ValidateJSON validates data against a JSON Schema
func ValidateJSON(schemaJSON []byte, data map[string]interface{}) error {
	var schema JSONSchema
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		return errors.New("invalid JSON schema format")
	}

	var errors []string
	validateJSONSchema(&schema, data, "", &errors)

	if len(errors) > 0 {
		return fmt.Errorf("schema validation failed:\n%s", errors)
	}
	return nil
}

// validateJSONSchema recursively validates data against the schema
func validateJSONSchema(schema *JSONSchema, data interface{}, path string, errors *[]string) {
	// Validate required fields
	if obj, ok := data.(map[string]interface{}); ok {
		for _, field := range schema.Required {
			if _, exists := obj[field]; !exists {
				*errors = append(*errors, fmt.Sprintf("%s: missing required field '%s'", path, field))
			}
		}
	}

	// Validate properties
	if schema.Properties != nil {
		obj, ok := data.(map[string]interface{})
		if !ok {
			*errors = append(*errors, fmt.Sprintf("%s: expected an object", path))
			return
		}
		for field, propSchema := range schema.Properties {
			if value, exists := obj[field]; exists {
				validateJSONSchema(propSchema, value, fmt.Sprintf("%s.%s", path, field), errors)
			}
		}
	}

	// Validate type(s)
	if schema.Type != nil {
		if !validateJSONType(schema.Type, data) {
			*errors = append(*errors, fmt.Sprintf("%s: incorrect type, expected %v", path, schema.Type))
		}
	}

	// Validate enum values
	if len(schema.Enum) > 0 {
		valid := false
		for _, e := range schema.Enum {
			if reflect.DeepEqual(e, data) {
				valid = true
				break
			}
		}
		if !valid {
			*errors = append(*errors, fmt.Sprintf("%s: value not in enum %v", path, schema.Enum))
		}
	}

	// Validate number constraints
	if num, ok := data.(float64); ok {
		if schema.Minimum != nil && num < *schema.Minimum {
			*errors = append(*errors, fmt.Sprintf("%s: value must be >= %f", path, *schema.Minimum))
		}
		if schema.Maximum != nil && num > *schema.Maximum {
			*errors = append(*errors, fmt.Sprintf("%s: value must be <= %f", path, *schema.Maximum))
		}
	}

	// Validate string constraints
	if str, ok := data.(string); ok {
		if schema.MinLength != nil && len(str) < *schema.MinLength {
			*errors = append(*errors, fmt.Sprintf("%s: string length must be >= %d", path, *schema.MinLength))
		}
		if schema.MaxLength != nil && len(str) > *schema.MaxLength {
			*errors = append(*errors, fmt.Sprintf("%s: string length must be <= %d", path, *schema.MaxLength))
		}
		if schema.Pattern != nil {
			if matched, _ := regexp.MatchString(*schema.Pattern, str); !matched {
				*errors = append(*errors, fmt.Sprintf("%s: does not match pattern '%s'", path, *schema.Pattern))
			}
		}
	}

	// Validate array constraints
	if arr, ok := data.([]interface{}); ok {
		if schema.MinItems != nil && len(arr) < *schema.MinItems {
			*errors = append(*errors, fmt.Sprintf("%s: array must have at least %d items", path, *schema.MinItems))
		}
		if schema.MaxItems != nil && len(arr) > *schema.MaxItems {
			*errors = append(*errors, fmt.Sprintf("%s: array must have at most %d items", path, *schema.MaxItems))
		}
		if schema.UniqueItems {
			uniqueSet := make(map[interface{}]bool)
			for _, v := range arr {
				if uniqueSet[v] {
					*errors = append(*errors, fmt.Sprintf("%s: array must have unique items", path))
					break
				}
				uniqueSet[v] = true
			}
		}
		for i, v := range arr {
			validateJSONSchema(schema.Items, v, fmt.Sprintf("%s[%d]", path, i), errors)
		}
	}
}

// validateJSONType checks if the value matches the expected JSON type(s)
func validateJSONType(expectedType interface{}, value interface{}) bool {
	switch expected := expectedType.(type) {
	case string:
		return checkJsonPrimitiveType(expected, value)
	case []interface{}: // Union type
		for _, subType := range expected {
			if validateJSONType(subType, value) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// checkJsonPrimitiveType validates JSON primitive types
func checkJsonPrimitiveType(expected string, value interface{}) bool {
	switch expected {
	case "null":
		return value == nil
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "integer":
		_, ok := value.(int)
		return ok
	case "number":
		_, ok := value.(float64)
		return ok
	case "string":
		_, ok := value.(string)
		return ok
	case "object":
		_, ok := value.(map[string]interface{})
		return ok
	case "array":
		_, ok := value.([]interface{})
		return ok
	default:
		return false
	}
}
