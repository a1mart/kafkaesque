package schemavalidator

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

// ValidateProto validates data against a protobuf-like schema (using raw data)
func ValidateProto(schemaContent []byte, data interface{}) error {
	var allErrors []string

	// Parse the expected fields from the schema
	var expectedFields map[string]FieldInfo
	err := json.Unmarshal(schemaContent, &expectedFields)
	if err != nil {
		return errors.New("failed to parse schema content: " + err.Error())
	}

	// Validate the data
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Map {
		return errors.New("expected map-like structure for proto validation")
	}

	// Iterate through expected fields and validate
	for fieldName, fieldInfo := range expectedFields {
		fieldVal := v.MapIndex(reflect.ValueOf(fieldName))

		// Check if field is missing
		if !fieldVal.IsValid() {
			if fieldInfo.Required {
				allErrors = append(allErrors, "missing required field: "+fieldName)
			}
			continue
		}

		// Handle repeated fields (arrays/slices)
		if fieldInfo.IsRepeated {
			if fieldVal.Kind() != reflect.Slice {
				allErrors = append(allErrors, "field "+fieldName+" should be a repeated field (slice)")
				continue
			}
			// Validate each element in the slice
			for i := 0; i < fieldVal.Len(); i++ {
				elem := fieldVal.Index(i)
				if elem.Type().String() != fieldInfo.ExpectedType {
					allErrors = append(allErrors, "element in field "+fieldName+" should be of type "+fieldInfo.ExpectedType)
				}
			}
		} else {
			// Validate field type
			if fieldVal.Kind() == reflect.Map {
				// Recursively validate nested fields (nested messages)
				nestedErrors := ValidateProto([]byte(`{"fields":`+fieldInfo.NestedFieldsToJSON()+`}`), fieldVal.Interface())
				if nestedErrors != nil {
					allErrors = append(allErrors, "nested validation errors for field "+fieldName+": "+nestedErrors.Error())
				}
			} else if fieldVal.Type().String() != fieldInfo.ExpectedType {
				allErrors = append(allErrors, "field "+fieldName+" should be of type "+fieldInfo.ExpectedType)
			}
		}
	}

	if len(allErrors) > 0 {
		return errors.New(strings.Join(allErrors, "; "))
	}
	return nil
}

// FieldInfo defines additional information about expected fields
type FieldInfo struct {
	ExpectedType string               // Expected field type (e.g., "string", "int32")
	Required     bool                 // Whether the field is required
	IsRepeated   bool                 // Whether the field is a repeated field (array/slice)
	NestedFields map[string]FieldInfo // Fields in nested messages
}

// NestedFieldsToJSON converts nested fields to JSON format
func (f *FieldInfo) NestedFieldsToJSON() string {
	nestedJSON, _ := json.Marshal(f.NestedFields)
	return string(nestedJSON)
}
