package validators

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Struct metadata cache
var structCache sync.Map // Stores map[reflect.Type]*structMetadata

type structMetadata struct {
	fields []fieldMetadata
}

type fieldMetadata struct {
	name      string
	validate  []string
	isStruct  bool
	isSlice   bool
	isPointer bool
}

// Main validation function
func Validate(val interface{}) *ValidationResult {
	sanitize(val) // Assuming this function mutates `val`
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem() // Dereference pointer
	}

	result := &ValidationResult{}
	validate(v, "", result)

	if result.HasErrors() {
		return result
	}
	return nil
}

func validate(v reflect.Value, fieldPath string, result *ValidationResult) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		result.AddError(fieldPath, "expected a struct")
		return
	}

	// Get cached struct metadata
	t := v.Type()
	meta, _ := structCache.Load(t)
	if meta == nil {
		meta = cacheStructMetadata(t)
	}

	fields := meta.(*structMetadata).fields

	for _, fieldMeta := range fields {
		field := v.FieldByName(fieldMeta.name)

		// Build field path efficiently
		fullPath := fieldMeta.name
		if fieldPath != "" {
			fullPath = fieldPath + "." + fieldMeta.name
		}

		// Dereference pointer fields
		if fieldMeta.isPointer && !field.IsNil() {
			field = field.Elem()
		}

		// Handle nested structs
		if fieldMeta.isStruct {
			validate(field, fullPath, result)
			continue
		}

		// Handle slices/arrays
		if fieldMeta.isSlice {
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.Ptr && !elem.IsNil() {
					elem = elem.Elem()
				}
				validate(elem, fullPath+"["+itoa(j)+"]", result)
			}
			continue
		}

		// Apply validation rules
		for _, rule := range fieldMeta.validate {
			if err := applyRule(rule, field, fieldMeta.name); err != nil {
				result.AddError(fullPath, err.Error())
			}
		}
	}
}

// Caches struct metadata to avoid reflection on every request
func cacheStructMetadata(t reflect.Type) *structMetadata {
	meta := &structMetadata{}
	numFields := t.NumField()
	meta.fields = make([]fieldMetadata, numFields)

	for i := 0; i < numFields; i++ {
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("validate")

		meta.fields[i] = fieldMetadata{
			name:      fieldType.Name,
			validate:  parseValidationTag(tag),
			isStruct:  fieldType.Type.Kind() == reflect.Struct,
			isSlice:   fieldType.Type.Kind() == reflect.Slice || fieldType.Type.Kind() == reflect.Array,
			isPointer: fieldType.Type.Kind() == reflect.Ptr,
		}
	}

	structCache.Store(t, meta)
	return meta
}

// Parses validation rules from the struct tag
func parseValidationTag(tag string) []string {
	if tag == "" {
		return nil
	}
	return strings.Split(tag, ",")
}

// Optimized integer to string conversion for small numbers
func itoa(n int) string {
	if n < 10 {
		return string('0' + byte(n)) // Fast conversion for small numbers
	}
	return strconv.Itoa(n) // Fallback for larger numbers
}
