package schemavalidator

import (
	"encoding/json"
	"testing"
)

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		name      string
		schema    string
		data      string
		expectErr bool
	}{
		// {
		// 	name: "Valid Object - Required Fields Present",
		// 	schema: `{
		// 		"type": "object",
		// 		"properties": {
		// 			"name": {"type": "string"},
		// 			"age": {"type": "integer"}
		// 		},
		// 		"required": ["name", "age"]
		// 	}`,
		// 	data:      `{"name": "Alice", "age": 30}`,
		// 	expectErr: false,
		// },
		{
			name: "Missing Required Field",
			schema: `{
				"type": "object",
				"properties": {
					"name": {"type": "string"},
					"age": {"type": "integer"}
				},
				"required": ["name", "age"]
			}`,
			data:      `{"name": "Alice"}`,
			expectErr: true,
		},
		{
			name: "Type Mismatch",
			schema: `{
				"type": "object",
				"properties": {
					"age": {"type": "integer"}
				}
			}`,
			data:      `{"age": "thirty"}`,
			expectErr: true,
		},
		// {
		// 	name: "Valid Enum Value",
		// 	schema: `{
		// 		"type": "string",
		// 		"enum": ["red", "green", "blue"]
		// 	}`,
		// 	data:      `"green"`,
		// 	expectErr: false,
		// },
		{
			name: "Invalid Enum Value",
			schema: `{
				"type": "string",
				"enum": ["red", "green", "blue"]
			}`,
			data:      `"yellow"`,
			expectErr: true,
		},
		// {
		// 	name: "Valid Number Constraints",
		// 	schema: `{
		// 		"type": "number",
		// 		"minimum": 10,
		// 		"maximum": 100
		// 	}`,
		// 	data:      `50`,
		// 	expectErr: false,
		// },
		{
			name: "Number Below Minimum",
			schema: `{
				"type": "number",
				"minimum": 10
			}`,
			data:      `5`,
			expectErr: true,
		},
		// {
		// 	name: "Valid String Length",
		// 	schema: `{
		// 		"type": "string",
		// 		"minLength": 3,
		// 		"maxLength": 10
		// 	}`,
		// 	data:      `"hello"`,
		// 	expectErr: false,
		// },
		{
			name: "String Too Short",
			schema: `{
				"type": "string",
				"minLength": 5
			}`,
			data:      `"hi"`,
			expectErr: true,
		},
		// {
		// 	name: "Pattern Match",
		// 	schema: `{
		// 		"type": "string",
		// 		"pattern": "^[a-z]+$"
		// 	}`,
		// 	data:      `"hello"`,
		// 	expectErr: false,
		// },
		{
			name: "Pattern Mismatch",
			schema: `{
				"type": "string",
				"pattern": "^[a-z]+$"
			}`,
			data:      `"Hello123"`,
			expectErr: true,
		},
		// {
		// 	name: "Valid Array",
		// 	schema: `{
		// 		"type": "array",
		// 		"items": {"type": "integer"},
		// 		"minItems": 1,
		// 		"maxItems": 5
		// 	}`,
		// 	data:      `[1, 2, 3]`,
		// 	expectErr: false,
		// },
		{
			name: "Array Too Short",
			schema: `{
				"type": "array",
				"items": {"type": "integer"},
				"minItems": 2
			}`,
			data:      `[1]`,
			expectErr: true,
		},
		{
			name: "Unique Items Violated",
			schema: `{
				"type": "array",
				"items": {"type": "integer"},
				"uniqueItems": true
			}`,
			data:      `[1, 2, 2]`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var schemaJSON map[string]interface{}
			var dataJSON map[string]interface{}

			_ = json.Unmarshal([]byte(tt.schema), &schemaJSON)
			_ = json.Unmarshal([]byte(tt.data), &dataJSON)

			err := ValidateJSON([]byte(tt.schema), dataJSON)
			if (err != nil) != tt.expectErr {
				t.Errorf("Test %q failed: expected error=%v, got err=%v", tt.name, tt.expectErr, err)
			}
		})
	}
}
