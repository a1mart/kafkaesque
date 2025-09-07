package schemavalidator

import (
	"testing"
)

// Test valid JSON against a correct Avro schema
func TestValidateAvro_ValidInput(t *testing.T) {
	schema := `{
		"type": "record",
		"name": "TestRecord",
		"fields": [
			{"name": "id", "type": "string"},
			{"name": "age", "type": "int"},
			{"name": "active", "type": "boolean"}
		]
	}`

	data := map[string]interface{}{
		"id":     "123",
		"age":    30,
		"active": true,
	}

	err := ValidateAvro([]byte(schema), data)
	if err != nil {
		t.Errorf("Expected valid input to pass, but got error: %v", err)
	}
}

// Test missing required fields
func TestValidateAvro_MissingFields(t *testing.T) {
	schema := `{
		"type": "record",
		"name": "TestRecord",
		"fields": [
			{"name": "id", "type": "string"},
			{"name": "age", "type": "int"}
		]
	}`

	data := map[string]interface{}{
		"id": "123",
	}

	err := ValidateAvro([]byte(schema), data)
	if err == nil {
		t.Errorf("Expected error due to missing required field, but got none")
	}
}

// Test incorrect data types
func TestValidateAvro_IncorrectTypes(t *testing.T) {
	schema := `{
		"type": "record",
		"name": "TestRecord",
		"fields": [
			{"name": "id", "type": "string"},
			{"name": "age", "type": "int"}
		]
	}`

	data := map[string]interface{}{
		"id":  123,  // Incorrect type (should be string)
		"age": "30", // Incorrect type (should be int)
	}

	err := ValidateAvro([]byte(schema), data)
	if err == nil {
		t.Errorf("Expected type mismatch error, but got none")
	}
}

// Test union types
func TestValidateAvro_UnionTypes(t *testing.T) {
	schema := `{
		"type": "record",
		"name": "TestRecord",
		"fields": [
			{"name": "value", "type": ["null", "string", "int"]}
		]
	}`

	validData1 := map[string]interface{}{"value": "hello"}
	validData2 := map[string]interface{}{"value": 42}
	validData3 := map[string]interface{}{"value": nil}
	invalidData := map[string]interface{}{"value": true} // Not in union

	if err := ValidateAvro([]byte(schema), validData1); err != nil {
		t.Errorf("Expected valid union input, but got error: %v", err)
	}
	if err := ValidateAvro([]byte(schema), validData2); err != nil {
		t.Errorf("Expected valid union input, but got error: %v", err)
	}
	if err := ValidateAvro([]byte(schema), validData3); err != nil {
		t.Errorf("Expected valid union input, but got error: %v", err)
	}
	if err := ValidateAvro([]byte(schema), invalidData); err == nil {
		t.Errorf("Expected error for invalid union type, but got none")
	}
}

// Test array type
func TestValidateAvro_ArrayType(t *testing.T) {
	schema := `{
		"type": "record",
		"name": "TestRecord",
		"fields": [
			{"name": "tags", "type": {"type": "array", "items": "string"}}
		]
	}`

	validData := map[string]interface{}{"tags": []interface{}{"tag1", "tag2"}}
	invalidData := map[string]interface{}{"tags": "not_an_array"}

	if err := ValidateAvro([]byte(schema), validData); err != nil {
		t.Errorf("Expected valid array input, but got error: %v", err)
	}
	if err := ValidateAvro([]byte(schema), invalidData); err == nil {
		t.Errorf("Expected error for non-array value, but got none")
	}
}

// Test map type
func TestValidateAvro_MapType(t *testing.T) {
	schema := `{
		"type": "record",
		"name": "TestRecord",
		"fields": [
			{"name": "metadata", "type": {"type": "map", "values": "string"}}
		]
	}`

	validData := map[string]interface{}{"metadata": map[string]interface{}{"key1": "value1", "key2": "value2"}}
	invalidData := map[string]interface{}{"metadata": "not_a_map"}

	if err := ValidateAvro([]byte(schema), validData); err != nil {
		t.Errorf("Expected valid map input, but got error: %v", err)
	}
	if err := ValidateAvro([]byte(schema), invalidData); err == nil {
		t.Errorf("Expected error for non-map value, but got none")
	}
}

// Test nested record
func TestValidateAvro_NestedRecord(t *testing.T) {
	schema := `{
		"type": "record",
		"name": "ParentRecord",
		"fields": [
			{"name": "child", "type": {
				"type": "record",
				"name": "ChildRecord",
				"fields": [
					{"name": "name", "type": "string"},
					{"name": "age", "type": "int"}
				]
			}}
		]
	}`

	validData := map[string]interface{}{
		"child": map[string]interface{}{
			"name": "Alice",
			"age":  10,
		},
	}
	invalidData := map[string]interface{}{
		"child": map[string]interface{}{
			"name": 10,    // Wrong type (should be string)
			"age":  "ten", // Wrong type (should be int)
		},
	}

	if err := ValidateAvro([]byte(schema), validData); err != nil {
		t.Errorf("Expected valid nested record, but got error: %v", err)
	}
	if err := ValidateAvro([]byte(schema), invalidData); err == nil {
		t.Errorf("Expected error for invalid nested record, but got none")
	}
}

// Test invalid schema format
func TestValidateAvro_InvalidSchema(t *testing.T) {
	invalidSchema := `{ "invalid_json": "missing_closing_brace"`

	data := map[string]interface{}{
		"field": "value",
	}

	err := ValidateAvro([]byte(invalidSchema), data)
	if err == nil {
		t.Errorf("Expected error for invalid schema format, but got none")
	}
}

// Test top-level non-record schema
func TestValidateAvro_NonRecordSchema(t *testing.T) {
	schema := `{
		"type": "string"
	}`

	data := map[string]interface{}{
		"field": "value",
	}

	err := ValidateAvro([]byte(schema), data)
	if err == nil {
		t.Errorf("Expected error for non-record schema, but got none")
	}
}
