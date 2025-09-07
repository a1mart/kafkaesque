package schemavalidator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AvroSchema represents an Avro schema
type AvroSchema struct {
	Type   interface{}       `json:"type"`
	Name   string            `json:"name,omitempty"`
	Fields []AvroSchemaField `json:"fields,omitempty"`
	Items  interface{}       `json:"items,omitempty"`  // For arrays
	Values interface{}       `json:"values,omitempty"` // For maps
}

// AvroSchemaField represents a field in an Avro record
type AvroSchemaField struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"`
}

// ValidateAvro validates JSON data against an Avro schema and returns all errors
func ValidateAvro(schemaJSON []byte, data map[string]interface{}) error {
	var schema AvroSchema
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		return fmt.Errorf("invalid Avro schema format: %v", err)
	}

	// Top-level schema must be a record
	if schema.Type != "record" {
		return fmt.Errorf("only 'record' type schemas are supported at the top level, found: %v", schema.Type)
	}

	var errors []string

	// Validate each field
	for _, field := range schema.Fields {
		val, exists := data[field.Name]
		if !exists {
			errors = append(errors, fmt.Sprintf("missing required field: %s", field.Name))
			continue
		}

		// Validate field type
		if err := validateType(field.Type, val, field.Name); err != nil {
			errors = append(errors, err.Error())
		}
	}

	// Return all collected errors
	if len(errors) > 0 {
		return fmt.Errorf("schema validation failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// validateType checks if a value matches the expected Avro type
func validateType(expectedType interface{}, value interface{}, fieldName string) error {
	switch t := expectedType.(type) {
	case string:
		// Handle primitive types correctly
		if !checkPrimitiveType(t, value) {
			return fmt.Errorf("field '%s' has incorrect type: expected %s, got %T", fieldName, t, value)
		}
	case []interface{}: // Union type
		for _, subType := range t {
			if validateType(subType, value, fieldName) == nil {
				return nil
			}
		}
		return fmt.Errorf("field '%s' does not match any allowed types: %v", fieldName, t)
	case map[string]interface{}:
		if typeName, ok := t["type"].(string); ok {
			switch typeName {
			case "array":
				return checkArrayType(t, value, fieldName)
			case "map":
				return checkMapType(t, value, fieldName)
			case "record":
				return checkRecordType(t, value, fieldName)
			default:
				return fmt.Errorf("field '%s' has unsupported complex type: %s", fieldName, typeName)
			}
		}
	default:
		return fmt.Errorf("field '%s' has unsupported type: %v", fieldName, expectedType)
	}
	return nil
}

// checkPrimitiveType validates Avro primitive types
func checkPrimitiveType(avroType string, value interface{}) bool {
	switch avroType {
	case "null":
		return value == nil
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "int":
		_, ok := value.(int) // Direct integer check
		if !ok {
			_, ok = value.(float64) // JSON unmarshals numbers as float64
			return ok
		}
		return true
	case "long":
		_, ok := value.(int64)
		if !ok {
			_, ok = value.(float64) // Handle JSON conversion issue
			return ok
		}
		return true
	case "float":
		_, ok := value.(float32)
		if !ok {
			_, ok = value.(float64) // Handle JSON conversion issue
			return ok
		}
		return true
	case "double":
		_, ok := value.(float64)
		return ok
	case "string":
		_, ok := value.(string)
		return ok
	case "bytes":
		_, ok := value.([]byte)
		return ok
	default:
		return false
	}
}

// checkArrayType validates an Avro array
func checkArrayType(schema map[string]interface{}, value interface{}, fieldName string) error {
	arr, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("field '%s' must be an array", fieldName)
	}

	itemsType := schema["items"]
	var errors []string
	for i, v := range arr {
		if err := validateType(itemsType, v, fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

// checkMapType validates an Avro map
func checkMapType(schema map[string]interface{}, value interface{}, fieldName string) error {
	m, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("field '%s' must be a map", fieldName)
	}

	valuesType := schema["values"]
	var errors []string
	for key, v := range m {
		if err := validateType(valuesType, v, fmt.Sprintf("%s[%s]", fieldName, key)); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

// checkRecordType validates an Avro record (nested struct)
func checkRecordType(schema map[string]interface{}, value interface{}, fieldName string) error {
	recordData, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("field '%s' must be a record (object)", fieldName)
	}

	fields, ok := schema["fields"].([]interface{})
	if !ok {
		return fmt.Errorf("field '%s' has an invalid record schema", fieldName)
	}

	var errors []string
	for _, field := range fields {
		fieldMap, ok := field.(map[string]interface{})
		if !ok {
			continue
		}

		fieldNameNested, ok := fieldMap["name"].(string)
		if !ok {
			continue
		}

		fieldType := fieldMap["type"]
		val, exists := recordData[fieldNameNested]
		if !exists {
			errors = append(errors, fmt.Sprintf("missing required field: %s.%s", fieldName, fieldNameNested))
			continue
		}

		if err := validateType(fieldType, val, fmt.Sprintf("%s.%s", fieldName, fieldNameNested)); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}
