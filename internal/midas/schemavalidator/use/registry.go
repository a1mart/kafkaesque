package main

import (
	"fmt"
	"log"

	"github.com/a1mart/kafkaesque/internal/midas/schemavalidator"
)

func main() {
	fmt.Println("*****DIRECT*****")
	//avro
	schema := `{
		"type": "record",
		"name": "User",
		"fields": [
			{"name": "id", "type": "int"},
			{"name": "name", "type": "string"},
			{"name": "email", "type": "string"}
		]
	}`

	data := map[string]interface{}{
		"id":    1,
		"name":  "Alice",
		"email": "alice@example.com",
	}

	err := schemavalidator.ValidateAvro([]byte(schema), data)
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation successful!")
	}

	//json
	schema = `{
		"type": "object",
		"properties": {
			"id": {"type": "integer"},
			"name": {"type": "string"},
			"email": {"type": "string"}
		},
		"required": ["id", "name", "email"]
	}`

	data = map[string]interface{}{
		"id":    1,
		"name":  "Alice",
		"email": "alice@example.com",
	}

	err = schemavalidator.ValidateJSON([]byte(schema), data)
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation successful!")
	}

	//protocol buffers
	// Define a schema for a protobuf-like data structure
	schema = `{
		"fields": {
			"Title": {"ExpectedType": "string", "Required": true},
			"Content": {"ExpectedType": "string", "Required": true},
			"Author": {"ExpectedType": "map", "Required": true, "NestedFields": {
				"Name": {"ExpectedType": "string", "Required": true},
				"Age": {"ExpectedType": "int32", "Required": true},
				"Email": {"ExpectedType": "string", "Required": true}
			}},
			"Tags": {"ExpectedType": "string", "IsRepeated": true}
		}
	}`

	// Define data to validate against the schema
	data = map[string]interface{}{
		"Title":   "My Post",
		"Content": "Content here...",
		"Author": map[string]interface{}{
			"Name":  "John",
			"Age":   30,
			"Email": "john@example.com",
		},
		"Tags": []interface{}{"Go", "Programming"},
	}

	// Validate the data
	err = schemavalidator.ValidateProto([]byte(schema), data)
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation successful!")
	}

	fmt.Println("*****CENTRAL REGISTRY*****")
	// Create schema registry
	registry := schemavalidator.NewSchemaRegistry()

	// Register schemas with various data types
	id, err := registry.RegisterSchema("", "ComplexExample", "v1", "json", `{
	"type": "object",
	"properties": {
		"id": {"type": "integer"},
		"name": {"type": "string"},
		"email": {"type": "string"},
		"active": {"type": "boolean"},
		"tags": {
			"type": "array",
			"items": {"type": "string"}
		},
		"metadata": {
			"type": "object",
			"properties": {
				"age": {"type": "integer"},
				"premium_user": {"type": "boolean"}
			},
			"required": ["age"]
		}
	},
	"required": ["id", "name", "email", "active", "metadata"]
}`)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the validator
	validator := schemavalidator.NewSchemaValidator(registry)

	// Example valid data
	validData := map[string]interface{}{
		"id":     123,
		"name":   "John Doe",
		"email":  "john@example.com",
		"active": true,
		"tags":   []interface{}{"golang", "backend", "microservices"},
		"metadata": map[string]interface{}{
			"age":          30,
			"premium_user": true,
		},
	}

	// Validate valid data using name/version
	err = validator.Validate("ComplexExample", "v1", validData)
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation successful! ✅")
	}

	// Example invalid data
	invalidData := map[string]interface{}{
		"id":    "wrong_type",
		"name":  456,
		"email": "invalid@example.com",
		"tags":  "not_an_array",
		"metadata": map[string]interface{}{
			"premium_user": "should_be_boolean",
		},
	}

	// Validate invalid data
	err = validator.Validate(id, "", invalidData)
	if err != nil {
		fmt.Println("Validation failed as expected: ❌", err)
	} else {
		fmt.Println("Unexpected success, check validation logic!")
	}

	// Register JSON schema
	_, err = registry.RegisterSchema("", "User", "v1", "json", `{
	"type": "object",
	"properties": {
		"id": {"type": "integer"},
		"name": {"type": "string"},
		"email": {"type": "string"}
	},
	"required": ["id", "name", "email"]
}`)
	if err != nil {
		log.Fatal(err)
	}

	// Register Avro schema
	avroId, err := registry.RegisterSchema("", "User", "v2", "avro", `{
	"type": "record",
	"name": "User",
	"fields": [
		{"name": "id", "type": "int"},
		{"name": "name", "type": "string"},
		{"name": "email", "type": "string"}
	]
}`)
	if err != nil {
		log.Fatal(err)
	}

	// Register Protocol Buffers schema
	_, err = registry.RegisterSchema("", "Post", "v1", "proto", `{
	"fields": {
		"Title": {"type": "string", "required": true},
		"Content": {"type": "string", "required": true},
		"Author": {
			"type": "User",
			"required": true,
			"nested": {
				"Name": {"type": "string", "required": true},
				"Age": {"type": "int32", "required": true},
				"Email": {"type": "string", "required": true}
			}
		},
		"Tags": {"type": "string", "repeated": true}
	}
}`)
	if err != nil {
		log.Fatal(err)
	}

	// Validate JSON schema
	data = map[string]interface{}{
		"id":    1,
		"name":  "Alice",
		"email": "alice@example.com",
	}
	err = validator.Validate("User", "v1", data)
	if err != nil {
		fmt.Println("JSON Validation failed:", err)
	} else {
		fmt.Println("JSON Validation successful!")
	}

	// Validate Avro schema
	err = validator.Validate(avroId, "", data)
	if err != nil {
		fmt.Println("Avro Validation failed:", err)
	} else {
		fmt.Println("Avro Validation successful!")
	}

	// Validate Proto schema
	protoData := map[string]interface{}{
		"Title":   "My Post",
		"Content": "Content here...",
		"Author": map[string]interface{}{
			"Name":  "John",
			"Age":   30,
			"Email": "john@example.com",
		},
		"Tags": []string{"Go", "Programming"},
	}
	err = validator.Validate("Post", "v1", protoData)
	if err != nil {
		fmt.Println("Proto Validation failed:", err)
	} else {
		fmt.Println("Proto Validation successful!")
	}

}
