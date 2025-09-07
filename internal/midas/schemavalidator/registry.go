package schemavalidator

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
)

// SchemaRegistry holds registered schemas
type SchemaRegistry struct {
	schemasByID       map[string]*Schema
	schemasByNameVers map[string]*Schema
	mu                sync.RWMutex
}

// Schema represents a schema definition
type Schema struct {
	//namespace, object name, version (automatic versioning)
	ID      string
	Name    string
	Version string
	Format  string // "json", "avro", "proto"
	Content string // Schema content
}

// NewSchemaRegistry initializes a new schema registry
func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{
		schemasByID:       make(map[string]*Schema),
		schemasByNameVers: make(map[string]*Schema),
	}
}

// generateSchemaID creates a unique hash-based ID for the schema
func generateSchemaID(name, version, format, content string) string {
	hash := sha256.Sum256([]byte(name + version + format + content))
	return hex.EncodeToString(hash[:])
}

// RegisterSchema registers a schema with an optional ID. If no ID is provided, one is generated.
func (r *SchemaRegistry) RegisterSchema(id, name, version, format, content string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if id == "" {
		id = generateSchemaID(name, version, format, content)
	}

	if _, exists := r.schemasByID[id]; exists {
		return "", errors.New("schema ID already exists")
	}

	schemaKey := name + ":" + version
	if _, exists := r.schemasByNameVers[schemaKey]; exists {
		return "", errors.New("schema name and version already registered")
	}

	schema := &Schema{ID: id, Name: name, Version: version, Format: format, Content: content}
	r.schemasByID[id] = schema
	r.schemasByNameVers[schemaKey] = schema

	return id, nil
}

// GetSchema retrieves a schema by either ID or name/version
func (r *SchemaRegistry) GetSchema(identifier, version string) (*Schema, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// If version is empty, assume lookup by ID
	if version == "" {
		if schema, exists := r.schemasByID[identifier]; exists {
			return schema, nil
		}
	} else {
		schemaKey := identifier + ":" + version
		if schema, exists := r.schemasByNameVers[schemaKey]; exists {
			return schema, nil
		}
	}

	return nil, errors.New("schema not found")
}

// SchemaValidator validates instances against schemas
type SchemaValidator struct {
	registry *SchemaRegistry
}

// NewSchemaValidator initializes a schema validator
func NewSchemaValidator(registry *SchemaRegistry) *SchemaValidator {
	return &SchemaValidator{registry: registry}
}

// Validate validates data against a registered schema, using either ID or name/version
func (v *SchemaValidator) Validate(identifier, version string, data interface{}) error {
	schema, err := v.registry.GetSchema(identifier, version)
	if err != nil {
		return err
	}

	switch schema.Format {
	case "json":
		return ValidateJSON([]byte(schema.Content), data.(map[string]interface{}))
	case "avro":
		return ValidateAvro([]byte(schema.Content), data.(map[string]interface{}))
	case "proto":
		return ValidateProto([]byte(schema.Content), data)
	default:
		return errors.New("unsupported schema format")
	}
}
