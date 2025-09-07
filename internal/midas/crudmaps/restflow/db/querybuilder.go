package main

import (
	"fmt"
	"strings"
)

// QueryBuilder is a generic SQL query builder supporting advanced operations.
type QueryBuilder[T any] struct {
	selectFields []string
	table        string
	joins        []string
	conditions   []string
	groupBy      []string
	orderBy      []string
	limit        int
	offset       int
	subqueries   []string
	transactions []string
}

// NewQueryBuilder creates a new query builder for a given table.
func NewQueryBuilder[T any](table string) *QueryBuilder[T] {
	return &QueryBuilder[T]{
		table: table,
	}
}

// Select specifies the columns to retrieve.
func (qb *QueryBuilder[T]) Select(fields ...string) *QueryBuilder[T] {
	qb.selectFields = append(qb.selectFields, fields...)
	return qb
}

// Join adds a JOIN clause.
func (qb *QueryBuilder[T]) Join(table, on string) *QueryBuilder[T] {
	qb.joins = append(qb.joins, fmt.Sprintf("JOIN %s ON %s", table, on))
	return qb
}

// Where adds a WHERE clause with optional arguments.
func (qb *QueryBuilder[T]) Where(condition string, args ...interface{}) *QueryBuilder[T] {
	qb.conditions = append(qb.conditions, fmt.Sprintf(condition, args...))
	return qb
}

// GroupBy adds a GROUP BY clause.
func (qb *QueryBuilder[T]) GroupBy(fields ...string) *QueryBuilder[T] {
	qb.groupBy = append(qb.groupBy, fields...)
	return qb
}

// OrderBy adds an ORDER BY clause.
func (qb *QueryBuilder[T]) OrderBy(fields ...string) *QueryBuilder[T] {
	qb.orderBy = append(qb.orderBy, fields...)
	return qb
}

// Limit sets a LIMIT clause.
func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	qb.limit = limit
	return qb
}

// Offset sets an OFFSET clause.
func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
	qb.offset = offset
	return qb
}

// Subquery adds a subquery.
func (qb *QueryBuilder[T]) Subquery(subquery string) *QueryBuilder[T] {
	qb.subqueries = append(qb.subqueries, subquery)
	return qb
}

// Transaction adds a transaction statement.
func (qb *QueryBuilder[T]) Transaction(statement string) *QueryBuilder[T] {
	qb.transactions = append(qb.transactions, statement)
	return qb
}

// Build constructs the SQL query.
func (qb *QueryBuilder[T]) Build() string {
	fields := "*"
	if len(qb.selectFields) > 0 {
		fields = strings.Join(qb.selectFields, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", fields, qb.table)

	if len(qb.joins) > 0 {
		query += " " + strings.Join(qb.joins, " ")
	}

	if len(qb.conditions) > 0 {
		query += " WHERE " + strings.Join(qb.conditions, " AND ")
	}

	if len(qb.groupBy) > 0 {
		query += " GROUP BY " + strings.Join(qb.groupBy, ", ")
	}

	if len(qb.orderBy) > 0 {
		query += " ORDER BY " + strings.Join(qb.orderBy, ", ")
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	if len(qb.subqueries) > 0 {
		query += " " + strings.Join(qb.subqueries, " ")
	}

	if len(qb.transactions) > 0 {
		query = strings.Join(qb.transactions, "; ") + "; " + query
	}

	return query
}

// ExampleUsage demonstrates building a complex SQL query.
func ExampleUsage() {
	qb := NewQueryBuilder[any]("users")
	query := qb.Select("id", "name").
		Join("orders", "users.id = orders.user_id").
		Where("users.active = %d", 1).
		GroupBy("users.id").
		OrderBy("users.name", "orders.created_at DESC").
		Limit(10).
		Offset(20).
		Subquery("(SELECT COUNT(*) FROM orders)").
		Transaction("BEGIN TRANSACTION").
		Build()

	fmt.Println(query)
}

func main() {
	ExampleUsage()
}
