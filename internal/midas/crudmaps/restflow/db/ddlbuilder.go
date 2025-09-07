package main

import (
	"fmt"
	"strings"
)

// AdvancedDDLBuilder is a generic SQL query builder for DDL operations.
type AdvancedDDLBuilder[T any] struct {
	table                      string
	createTableStmt            string
	alterTableStmts            []string
	dropTableStmt              string
	createViewStmt             string
	createMaterializedViewStmt string
	dropViewStmt               string
	createIndexStmt            string
	dropIndexStmt              string
	createTriggerStmt          string
	dropTriggerStmt            string
	createFunctionStmt         string
	createProcedureStmt        string
	dropFunctionStmt           string
	dropProcedureStmt          string
	createSequenceStmt         string
	alterSequenceStmt          string
	dropSequenceStmt           string
	createDomainStmt           string
	createTypeStmt             string
	dropDomainStmt             string
	dropTypeStmt               string
}

// NewAdvancedDDLBuilder creates a new DDL builder for a given table.
func NewAdvancedDDLBuilder[T any](table string) *AdvancedDDLBuilder[T] {
	return &AdvancedDDLBuilder[T]{table: table}
}

// Column represents a column definition with type and optional constraints.
type Column struct {
	Name        string
	Type        string
	Constraints []string
}

// CreateTable generates the CREATE TABLE SQL statement with enhanced options.
func (adb *AdvancedDDLBuilder[T]) CreateTable(tableName string, columns []Column) *AdvancedDDLBuilder[T] {
	var columnDefs []string
	for _, col := range columns {
		colDef := fmt.Sprintf("%s %s", col.Name, col.Type)
		if len(col.Constraints) > 0 {
			colDef += " " + strings.Join(col.Constraints, " ")
		}
		columnDefs = append(columnDefs, colDef)
	}
	adb.createTableStmt = fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columnDefs, ", "))
	return adb
}

// AlterTable adds an ALTER TABLE SQL statement.
func (adb *AdvancedDDLBuilder[T]) AlterTable(alterations ...string) *AdvancedDDLBuilder[T] {
	adb.alterTableStmts = append(adb.alterTableStmts, alterations...)
	return adb
}

// DropTable generates the DROP TABLE SQL statement.
func (adb *AdvancedDDLBuilder[T]) DropTable(tableName string) *AdvancedDDLBuilder[T] {
	adb.dropTableStmt = fmt.Sprintf("DROP TABLE %s", tableName)
	return adb
}

// AddColumn adds a new column to an existing table.
func (adb *AdvancedDDLBuilder[T]) AddColumn(column Column) *AdvancedDDLBuilder[T] {
	alter := fmt.Sprintf("ADD COLUMN %s %s", column.Name, column.Type)
	if len(column.Constraints) > 0 {
		alter += " " + strings.Join(column.Constraints, " ")
	}
	adb.alterTableStmts = append(adb.alterTableStmts, alter)
	return adb
}

// DropColumn drops a column from an existing table.
func (adb *AdvancedDDLBuilder[T]) DropColumn(columnName string) *AdvancedDDLBuilder[T] {
	alter := fmt.Sprintf("DROP COLUMN %s", columnName)
	adb.alterTableStmts = append(adb.alterTableStmts, alter)
	return adb
}

// CreateForeignKey generates the CREATE FOREIGN KEY SQL statement.
func (adb *AdvancedDDLBuilder[T]) CreateForeignKey(constraintName, column, referencedTable, referencedColumn string) *AdvancedDDLBuilder[T] {
	alter := fmt.Sprintf("ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)", constraintName, column, referencedTable, referencedColumn)
	adb.alterTableStmts = append(adb.alterTableStmts, alter)
	return adb
}

// CreateIndex generates the CREATE INDEX SQL statement with options for unique and partial indexes.
func (adb *AdvancedDDLBuilder[T]) CreateIndex(indexName, tableName, columns string, unique bool, where string) *AdvancedDDLBuilder[T] {
	indexStmt := fmt.Sprintf("CREATE INDEX %s ON %s (%s)", indexName, tableName, columns)
	if unique {
		indexStmt = fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s (%s)", indexName, tableName, columns)
	}
	if where != "" {
		indexStmt += fmt.Sprintf(" WHERE %s", where)
	}
	adb.createIndexStmt = indexStmt
	return adb
}

// DropIndex generates the DROP INDEX SQL statement.
func (adb *AdvancedDDLBuilder[T]) DropIndex(indexName string) *AdvancedDDLBuilder[T] {
	adb.dropIndexStmt = fmt.Sprintf("DROP INDEX %s", indexName)
	return adb
}

// CreateTrigger generates the CREATE TRIGGER SQL statement with timing, event, and action.
func (adb *AdvancedDDLBuilder[T]) CreateTrigger(triggerName, tableName, timing, event, action string) *AdvancedDDLBuilder[T] {
	adb.createTriggerStmt = fmt.Sprintf("CREATE TRIGGER %s %s %s ON %s FOR EACH ROW EXECUTE PROCEDURE %s()", triggerName, timing, event, tableName, action)
	return adb
}

// DropTrigger generates the DROP TRIGGER SQL statement.
func (adb *AdvancedDDLBuilder[T]) DropTrigger(triggerName, tableName string) *AdvancedDDLBuilder[T] {
	adb.dropTriggerStmt = fmt.Sprintf("DROP TRIGGER %s ON %s", triggerName, tableName)
	return adb
}

// CreateFunction generates the CREATE FUNCTION SQL statement with parameters and return type.
func (adb *AdvancedDDLBuilder[T]) CreateFunction(functionName, returnType, body string, params map[string]string) *AdvancedDDLBuilder[T] {
	var paramDefs []string
	for param, paramType := range params {
		paramDefs = append(paramDefs, fmt.Sprintf("%s %s", param, paramType))
	}
	paramStr := strings.Join(paramDefs, ", ")
	adb.createFunctionStmt = fmt.Sprintf("CREATE FUNCTION %s(%s) RETURNS %s AS $$ %s $$ LANGUAGE plpgsql", functionName, paramStr, returnType, body)
	return adb
}

// CreateProcedure generates the CREATE PROCEDURE SQL statement.
func (adb *AdvancedDDLBuilder[T]) CreateProcedure(procedureName, body string, params map[string]string) *AdvancedDDLBuilder[T] {
	var paramDefs []string
	for param, paramType := range params {
		paramDefs = append(paramDefs, fmt.Sprintf("%s %s", param, paramType))
	}
	paramStr := strings.Join(paramDefs, ", ")
	adb.createProcedureStmt = fmt.Sprintf("CREATE PROCEDURE %s(%s) LANGUAGE plpgsql AS $$ %s $$", procedureName, paramStr, body)
	return adb
}

// DropFunction generates the DROP FUNCTION SQL statement.
func (adb *AdvancedDDLBuilder[T]) DropFunction(functionName string) *AdvancedDDLBuilder[T] {
	adb.dropFunctionStmt = fmt.Sprintf("DROP FUNCTION %s", functionName)
	return adb
}

// DropProcedure generates the DROP PROCEDURE SQL statement.
func (adb *AdvancedDDLBuilder[T]) DropProcedure(procedureName string) *AdvancedDDLBuilder[T] {
	adb.dropProcedureStmt = fmt.Sprintf("DROP PROCEDURE %s", procedureName)
	return adb
}

// CreateSequence generates the CREATE SEQUENCE SQL statement with options.
func (adb *AdvancedDDLBuilder[T]) CreateSequence(sequenceName, startValue, increment string) *AdvancedDDLBuilder[T] {
	adb.createSequenceStmt = fmt.Sprintf("CREATE SEQUENCE %s START %s INCREMENT %s", sequenceName, startValue, increment)
	return adb
}

// AlterSequence generates the ALTER SEQUENCE SQL statement.
func (adb *AdvancedDDLBuilder[T]) AlterSequence(sequenceName, alteration string) *AdvancedDDLBuilder[T] {
	adb.alterSequenceStmt = fmt.Sprintf("ALTER SEQUENCE %s %s", sequenceName, alteration)
	return adb
}

// DropSequence generates the DROP SEQUENCE SQL statement.
func (adb *AdvancedDDLBuilder[T]) DropSequence(sequenceName string) *AdvancedDDLBuilder[T] {
	adb.dropSequenceStmt = fmt.Sprintf("DROP SEQUENCE %s", sequenceName)
	return adb
}

// CreateDomain generates the CREATE DOMAIN SQL statement.
func (adb *AdvancedDDLBuilder[T]) CreateDomain(domainName, domainType string, checkClause string) *AdvancedDDLBuilder[T] {
	adb.createDomainStmt = fmt.Sprintf("CREATE DOMAIN %s AS %s", domainName, domainType)
	if checkClause != "" {
		adb.createDomainStmt += fmt.Sprintf(" CHECK (%s)", checkClause)
	}
	return adb
}

// CreateType generates the CREATE TYPE SQL statement.
func (adb *AdvancedDDLBuilder[T]) CreateType(typeName, baseType string) *AdvancedDDLBuilder[T] {
	adb.createTypeStmt = fmt.Sprintf("CREATE TYPE %s AS %s", typeName, baseType)
	return adb
}

// DropDomain generates the DROP DOMAIN SQL statement.
func (adb *AdvancedDDLBuilder[T]) DropDomain(domainName string) *AdvancedDDLBuilder[T] {
	adb.dropDomainStmt = fmt.Sprintf("DROP DOMAIN %s", domainName)
	return adb
}

// DropType generates the DROP TYPE SQL statement.
func (adb *AdvancedDDLBuilder[T]) DropType(typeName string) *AdvancedDDLBuilder[T] {
	adb.dropTypeStmt = fmt.Sprintf("DROP TYPE %s", typeName)
	return adb
}

// Build constructs all DDL statements.
func (adb *AdvancedDDLBuilder[T]) Build() string {
	var statements []string

	if adb.createTableStmt != "" {
		statements = append(statements, adb.createTableStmt)
	}
	if len(adb.alterTableStmts) > 0 {
		statements = append(statements, strings.Join(adb.alterTableStmts, "; "))
	}
	if adb.dropTableStmt != "" {
		statements = append(statements, adb.dropTableStmt)
	}
	if adb.createViewStmt != "" {
		statements = append(statements, adb.createViewStmt)
	}
	if adb.createMaterializedViewStmt != "" {
		statements = append(statements, adb.createMaterializedViewStmt)
	}
	if adb.dropViewStmt != "" {
		statements = append(statements, adb.dropViewStmt)
	}
	if adb.createIndexStmt != "" {
		statements = append(statements, adb.createIndexStmt)
	}
	if adb.dropIndexStmt != "" {
		statements = append(statements, adb.dropIndexStmt)
	}
	if adb.createTriggerStmt != "" {
		statements = append(statements, adb.createTriggerStmt)
	}
	if adb.dropTriggerStmt != "" {
		statements = append(statements, adb.dropTriggerStmt)
	}
	if adb.createFunctionStmt != "" {
		statements = append(statements, adb.createFunctionStmt)
	}
	if adb.createProcedureStmt != "" {
		statements = append(statements, adb.createProcedureStmt)
	}
	if adb.dropFunctionStmt != "" {
		statements = append(statements, adb.dropFunctionStmt)
	}
	if adb.dropProcedureStmt != "" {
		statements = append(statements, adb.dropProcedureStmt)
	}
	if adb.createSequenceStmt != "" {
		statements = append(statements, adb.createSequenceStmt)
	}
	if adb.alterSequenceStmt != "" {
		statements = append(statements, adb.alterSequenceStmt)
	}
	if adb.dropSequenceStmt != "" {
		statements = append(statements, adb.dropSequenceStmt)
	}
	if adb.createDomainStmt != "" {
		statements = append(statements, adb.createDomainStmt)
	}
	if adb.createTypeStmt != "" {
		statements = append(statements, adb.createTypeStmt)
	}
	if adb.dropDomainStmt != "" {
		statements = append(statements, adb.dropDomainStmt)
	}
	if adb.dropTypeStmt != "" {
		statements = append(statements, adb.dropTypeStmt)
	}

	return strings.Join(statements, "; ")
}

func main() {
	// Example usage
	builder := NewAdvancedDDLBuilder[string]("users")
	columns := []Column{
		{"id", "SERIAL", []string{"PRIMARY KEY"}},
		{"username", "VARCHAR(100)", []string{"NOT NULL", "UNIQUE"}},
		{"email", "VARCHAR(100)", []string{"NOT NULL", "UNIQUE"}},
		{"created_at", "TIMESTAMP", []string{"DEFAULT CURRENT_TIMESTAMP"}},
	}

	migration := builder.CreateTable("users", columns).
		CreateIndex("users_username_index", "users", "username", true, "").
		AlterTable("ADD COLUMN address VARCHAR(255)").
		Build()

	fmt.Println(migration)
}
