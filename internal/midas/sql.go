package validators

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// create table from structs

var (
	Struct2Table = map[string]string{
		"User": "id_users",
	}
)

// Struct2SQL and SQL2Struct

func GetSQLFieldMappings(val interface{}) map[string]string {
	mappings := make(map[string]string)
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		sqlTag := field.Tag.Get("sql")
		if sqlTag != "" {
			mappings[field.Name] = sqlTag
		}
	}
	return mappings
}

// GetSQLFieldMappingsForAction returns a map of field names to SQL column names for a specific action
func GetSQLFieldMappingsForAction(val interface{}, action string) map[string]string {
	mappings := make(map[string]string)
	v := reflect.ValueOf(val)
	fmt.Println(v, v.Type())
	// Check if the provided value is valid and not nil
	if !v.IsValid() {
		fmt.Println("Invalid value provided to GetSQLFieldMappingsForAction")
		return mappings
	}

	// Dereference if the value is a pointer
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			fmt.Println("Nil pointer provided to GetSQLFieldMappingsForAction")
			return mappings
		}
		v = v.Elem()
	}

	// Ensure the value is a struct
	if v.Kind() != reflect.Struct {
		fmt.Println("Non-struct type provided to GetSQLFieldMappingsForAction")
		return mappings
	}

	t := v.Type()

	// Iterate over struct fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		sqlTag := field.Tag.Get("sql") // Get the SQL tag value
		if sqlTag != "" {
			// Split the sql tag value (e.g., "username,select,insert")
			tags := strings.Split(sqlTag, ",")
			columnName := tags[0] // The first part is the column name
			actions := tags[1:]   // The rest are the allowed actions (e.g., select, insert, update, delete)

			// Check if the field is allowed for the given action
			if contains(actions, action) {
				mappings[columnName] = field.Name
			}
		}
	}
	return mappings
}

// GetSQLQueryForAction generates an SQL query based on the action and the struct's sql tags
func GetSQLQueryForAction(val interface{}, action string) (string, error) {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	// Initialize query parts
	var columns []string
	var values []string
	var setFields []string
	var whereConditions []string

	// Iterate over struct fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		sqlTag := field.Tag.Get("sql") // Get the sql tag
		if sqlTag != "" {
			tags := strings.Split(sqlTag, ",")
			columnName := tags[0] // The first part is the column name
			actions := tags[1:]   // The rest are actions

			// Depending on the action, add to appropriate query part
			if contains(actions, action) {
				//fieldName := field.Name
				fieldValue := fmt.Sprintf("'%v'", v.Field(i).Interface()) // Get the field value

				switch action {
				case "select":
					columns = append(columns, columnName)
				case "insert":
					columns = append(columns, columnName)
					values = append(values, fieldValue)
				case "update":
					setFields = append(setFields, fmt.Sprintf("%s = %s", columnName, fieldValue))
				case "delete":
					if columnName == "id" {
						whereConditions = append(whereConditions, fmt.Sprintf("%s = %s", columnName, fieldValue))
					}
				}
			}
		}
	}

	// Build the query based on action
	var query string
	switch action {
	case "select":
		query = fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), Struct2Table[t.Name()])
	case "insert":
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", Struct2Table[t.Name()], strings.Join(columns, ", "), strings.Join(values, ", "))
	case "update":
		query = fmt.Sprintf("UPDATE %s SET %s WHERE %s", Struct2Table[t.Name()], strings.Join(setFields, ", "), strings.Join(whereConditions, " AND "))
	case "delete":
		query = fmt.Sprintf("DELETE FROM %s WHERE %s", Struct2Table[t.Name()], strings.Join(whereConditions, " AND "))
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}

	return query, nil
}

// SQLQueryOptions holds all optional parameters for building an SQL query
type SQLQueryOptions struct {
	SelectFields   []string
	Filter         map[string]interface{}
	JoinConditions []string
	OrderBy        string
	GroupBy        string
	Limit          int
	Offset         int
}

// AdvancedSQL generates an SQL query based on the action, struct, and options
func AdvancedSQL(val interface{}, action string, options SQLQueryOptions) (string, error) {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	// Initialize query parts
	var columns []string
	var values []string
	var setFields []string
	var whereConditions []string
	var joinClauses []string

	// Handle struct fields and generate query components
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		sqlTag := field.Tag.Get("sql") // Get the sql tag
		if sqlTag != "" {
			tags := strings.Split(sqlTag, ",")
			columnName := tags[0]
			actions := tags[1:]

			if contains(actions, action) {
				fieldValue := v.Field(i)

				// Check if the field value is valid or nil
				var formattedValue string
				if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Interface {
					if fieldValue.IsNil() {
						formattedValue = "NULL" // Handle nil pointers or interfaces
					} else if fieldValue.Elem().Kind() == reflect.Struct && fieldValue.Type() == reflect.TypeOf(time.Time{}) {
						// Check if time.Time is zero
						t := fieldValue.Interface().(time.Time)
						if t.IsZero() {
							formattedValue = "NULL"
						} else {
							formattedValue = fmt.Sprintf("'%s'", t.Format("2006-01-02 15:04:05.999999-07:00"))
						}
					} else {
						formattedValue = fmt.Sprintf("'%v'", fieldValue.Interface())
					}
				} else if fieldValue.Kind() == reflect.Struct && fieldValue.Type() == reflect.TypeOf(time.Time{}) {
					// Direct time.Time handling
					t := fieldValue.Interface().(time.Time)
					if t.IsZero() {
						formattedValue = "NULL"
					} else {
						formattedValue = fmt.Sprintf("'%s'", t.Format("2006-01-02 15:04:05.999999-07:00"))
					}
				} else {
					formattedValue = fmt.Sprintf("'%v'", fieldValue.Interface()) // Default for non-nil values
				}

				switch action {
				case "select":
					columns = append(columns, columnName)
				case "insert":
					columns = append(columns, columnName)
					values = append(values, formattedValue)
				case "update":
					setFields = append(setFields, fmt.Sprintf("%s = %s", columnName, formattedValue))
				case "delete":
					if columnName == "id" {
						whereConditions = append(whereConditions, fmt.Sprintf("%s = %s", columnName, formattedValue))
					}
				}
			}
		}
	}

	// Handle additional filters for WHERE clause
	if len(options.Filter) > 0 {
		for fieldName, value := range options.Filter {
			whereConditions = append(whereConditions, fmt.Sprintf("%s = '%v'", fieldName, value))
		}
	}

	// Handle JOIN conditions
	if len(options.JoinConditions) > 0 {
		joinClauses = append(joinClauses, options.JoinConditions...)
	}

	// Build the query based on action
	var query string
	switch action {
	case "select":
		// Default to all columns if no specific ones are selected
		if len(columns) == 0 {
			columns = append(columns, "*")
		}
		query = fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), Struct2Table[t.Name()])

		// Add JOIN clauses
		if len(joinClauses) > 0 {
			query += " " + strings.Join(joinClauses, " ")
		}

		// Add WHERE conditions
		if len(whereConditions) > 0 {
			query += " WHERE " + strings.Join(whereConditions, " AND ")
		}

		// Add GROUP BY
		if options.GroupBy != "" {
			query += " GROUP BY " + options.GroupBy
		}

		// Add ORDER BY
		if options.OrderBy != "" {
			query += " ORDER BY " + options.OrderBy
		}

		// Add LIMIT and OFFSET
		if options.Limit > 0 {
			query += fmt.Sprintf(" LIMIT %d", options.Limit)
		}
		if options.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", options.Offset)
		}

	case "insert":
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", Struct2Table[t.Name()], strings.Join(columns, ", "), strings.Join(values, ", "))
	case "update":
		query = fmt.Sprintf("UPDATE %s SET %s WHERE %s", Struct2Table[t.Name()], strings.Join(setFields, ", "), strings.Join(whereConditions, " AND "))
	case "delete":
		query = fmt.Sprintf("DELETE FROM %s WHERE %s", Struct2Table[t.Name()], strings.Join(whereConditions, " AND "))
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}

	return query, nil
}
