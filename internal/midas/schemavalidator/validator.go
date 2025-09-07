package schemavalidator

// import "errors"

// // SchemaValidator validates instances against schemas
// type SchemaValidator struct {
// 	registry *SchemaRegistry
// }

// func NewSchemaValidator(registry *SchemaRegistry) *SchemaValidator {
// 	return &SchemaValidator{registry: registry}
// }

// // Validate validates data against a registered schema
// func (v *SchemaValidator) Validate(schemaName, version string, data interface{}) error {
// 	schema, err := v.registry.GetSchema(schemaName, version)
// 	if err != nil {
// 		return err
// 	}

// 	switch schema.Format {
// 	case "json":
// 		return ValidateJSON([]byte(schema.Content), data.(map[string]interface{}))
// 	case "avro":
// 		return ValidateAvro([]byte(schema.Content), data.(map[string]interface{}))
// 	case "proto":
// 		return ValidateProto([]byte(schema.Content), data)
// 	default:
// 		return errors.New("unsupported schema format")
// 	}
// }
