package validators

import (
	"fmt"
	"reflect"
	"strings"
)

func CheckFieldRoles(val interface{}, userRole string) error {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		roles := strings.Split(field.Tag.Get("role"), ",")
		if len(roles) > 0 && !contains(roles, userRole) {
			return fmt.Errorf("role %s is not authorized to access field %s", userRole, field.Name)
		}
	}
	return nil
}

func GetAllowedFieldsForRole(val interface{}, userRole string) ([]string, error) {
	var allowedFields []string

	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	// Iterate through the struct fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		roles := strings.Split(field.Tag.Get("role"), ",")

		// Check if the field contains the userRole in its "role" tag
		if len(roles) > 0 && contains(roles, userRole) {
			allowedFields = append(allowedFields, field.Name)
		}
	}

	if len(allowedFields) == 0 {
		return nil, fmt.Errorf("role %s is not authorized to access any fields", userRole)
	}

	return allowedFields, nil
}

func GetAllowedFieldsForRoleWithValues(val interface{}, userRole string) (map[string]interface{}, error) {
	allowedFields := make(map[string]interface{})

	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	// Iterate through the struct fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		roles := strings.Split(field.Tag.Get("role"), ",")

		// Check if the field contains the userRole in its "role" tag
		if len(roles) > 0 && contains(roles, userRole) {
			fieldValue := v.Field(i).Interface() // Get the field's value
			allowedFields[field.Name] = fieldValue
		}
	}

	if len(allowedFields) == 0 {
		return nil, fmt.Errorf("role %s is not authorized to access any fields", userRole)
	}

	return allowedFields, nil
}

func contains(slice []string, item string) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}
