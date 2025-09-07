package validators

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// Cache struct field sanitization metadata
var sanitizeCache sync.Map // map[reflect.Type]*structSanitizeMetadata

type structSanitizeMetadata struct {
	fields []sanitizeFieldMetadata
}

type sanitizeFieldMetadata struct {
	index int
	rules []string
}

// Main sanitization function
func sanitize(val interface{}) {
	v := reflect.ValueOf(val)

	// Handle pointers to structs
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	// Ensure the value is a struct
	if v.Kind() != reflect.Struct {
		fmt.Println("Sanitize can only be applied to structs")
		return
	}

	// Get cached metadata or compute it
	t := v.Type()
	meta, _ := sanitizeCache.Load(t)
	if meta == nil {
		meta = cacheSanitizeMetadata(t)
	}
	fields := meta.(*structSanitizeMetadata).fields

	// Iterate through fields using cached metadata
	for _, fieldMeta := range fields {
		field := v.Field(fieldMeta.index)

		// Skip unaddressable fields
		if !field.CanAddr() {
			continue
		}

		applySanitization(field.Addr(), fieldMeta.rules)
	}
}

// Caches struct metadata for sanitization
func cacheSanitizeMetadata(t reflect.Type) *structSanitizeMetadata {
	meta := &structSanitizeMetadata{}
	numFields := t.NumField()

	for i := 0; i < numFields; i++ {
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("sanitize")

		if tag == "" {
			continue
		}

		meta.fields = append(meta.fields, sanitizeFieldMetadata{
			index: i,
			rules: parseSanitizeTag(tag),
		})
	}

	sanitizeCache.Store(t, meta)
	return meta
}

// Parses sanitization rules from the struct tag
func parseSanitizeTag(tag string) []string {
	return strings.Split(tag, ",")
}

// Applies sanitization rules
func applySanitization(field reflect.Value, rules []string) {
	// Ensure field is addressable
	if field.Kind() == reflect.Ptr && !field.IsNil() {
		field = field.Elem()
	}

	// Apply rules
	for _, rule := range rules {
		switch rule {
		case "lowercase":
			if field.Kind() == reflect.String {
				field.SetString(strings.ToLower(field.String()))
			}
		case "trim":
			if field.Kind() == reflect.String {
				field.SetString(strings.TrimSpace(field.String()))
			}
		case "uppercase":
			if field.Kind() == reflect.String {
				field.SetString(strings.ToUpper(field.String()))
			}
		case "titlecase":
			if field.Kind() == reflect.String {
				s := field.String()
				if s != "" {
					field.SetString(strings.Title(strings.ToLower(s)))
				}
			}
		default:
			// Handle dynamic replace rule: "replace=old:new"
			if strings.HasPrefix(rule, "replace=") {
				replacePair := strings.SplitN(strings.TrimPrefix(rule, "replace="), ":", 2)
				if len(replacePair) == 2 && field.Kind() == reflect.String {
					field.SetString(strings.ReplaceAll(field.String(), replacePair[0], replacePair[1]))
				}
			} else {
				fmt.Printf("Unknown sanitization rule: %s\n", rule)
			}
		}
	}
}
