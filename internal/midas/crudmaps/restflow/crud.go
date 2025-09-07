package main

import (
	"fmt"
	"reflect"
	"strings"
)

/*
[Transport Security] -> [Context & Tracing] -> [Rate Limiting] -> [AuthN/AuthZ]
-> [Idempotency Checks] -> [Deserialize] -> [Sanitize] -> [Validate]
-> [Transform & Enrich] -> [External Services] -> [Transactions] -> [Storage]
-> [Cache] -> [Event Emission] -> [Search Indexing] -> [Encode Response]
-> [Post-Response Hooks] -> [Observability] -> [Backup/Recovery] -> [Shutdown]
-> [Config Management] -> [Schema Evolution]

CREATE (POST)
[Transport Security] -> [Context & Tracing] -> [Rate Limiting] -> [AuthN/AuthZ]
-> [Idempotency Checks] -> [Deserialize] -> [Sanitize] -> [Validate]
-> [Transform & Enrich] -> [External Services] -> [Transactions]
-> [Storage: Insert] -> [Cache: Invalidate/Set] -> [Event Emission]
-> [Search Indexing: Add] -> [Encode Response] -> [Post-Response Hooks]
-> [Observability] -> [Backup/Recovery] -> [Shutdown] -> [Config Management]
-> [Schema Evolution]
*Idempotency Checks: Important here to avoid duplicate resource creation.
*Validation: Strong focus on uniqueness and required fields.
*Transactions: Ensure DB changes only commit if everything succeeds.
*Event Emission: Notify other services a new entity exists.
*Search Indexing: Add the new resource.
*Cache: Prime cache for faster future reads.

READ (GET)
[Transport Security] -> [Context & Tracing] -> [Rate Limiting] -> [AuthN/AuthZ]
-> [Sanitize] -> [Validate] -> [Cache: Read]
-> [Storage: Select if cache miss] -> [Search Indexing: Query if needed]
-> [Encode Response] -> [Post-Response Hooks] -> [Observability]
-> [Shutdown] -> [Config Management] -> [Schema Evolution]
*Idempotency Checks: Not needed — GETs are idempotent by definition.
*Cache: Attempt cache hit first; only go to storage on miss.
*Transactions: Not needed — read-only operation.
*Search Indexing: Query if advanced search is needed.
*Event Emission: Not applicable — reads don’t generate events.

UDPATE (PUT/PATCH)
[Transport Security] -> [Context & Tracing] -> [Rate Limiting] -> [AuthN/AuthZ]
-> [Idempotency Checks] -> [Deserialize] -> [Sanitize] -> [Validate]
-> [Transform & Enrich] -> [External Services] -> [Transactions]
-> [Storage: Update] -> [Cache: Invalidate/Set] -> [Event Emission]
-> [Search Indexing: Update] -> [Encode Response] -> [Post-Response Hooks]
-> [Observability] -> [Backup/Recovery] -> [Shutdown] -> [Config Management]
-> [Schema Evolution]
*Idempotency Checks: Useful for PATCH, especially with retryable requests.
*Validation: May require checking existing state and cross-field consistency.
*Transactions: Ensure updates apply atomically.
*Cache: Invalidate old cache or update it with the new state.
*Event Emission: Notify listeners about changes.
*Search Indexing: Update the indexed document if the resource changes.

DELETE (DELETE)
[Transport Security] -> [Context & Tracing] -> [Rate Limiting] -> [AuthN/AuthZ]
-> [Idempotency Checks] -> [Sanitize] -> [Validate]
-> [External Services: Cleanup] -> [Transactions] -> [Storage: Delete]
-> [Cache: Invalidate] -> [Event Emission] -> [Search Indexing: Remove]
-> [Encode Response] -> [Post-Response Hooks] -> [Observability]
-> [Backup/Recovery] -> [Shutdown] -> [Config Management] -> [Schema Evolution]
*Idempotency Checks: Can avoid repeated deletion attempts leading to errors.
*Validation: Confirm the resource exists and can be safely removed.
*Transactions: Ensure cascading deletions and side-effects are consistent.
*Cache: Ensure stale data is cleared.
*Event Emission: Notify other services the resource no longer exists.
*Search Indexing: Remove from search engine index.

BULK OPERATIONS (BATCH CREATE/UPDATE/DELETE)
[Transport Security] -> [Context & Tracing] -> [Rate Limiting] -> [AuthN/AuthZ]
-> [Idempotency Checks] -> [Deserialize] -> [Sanitize] -> [Validate]
-> [Transform & Enrich] -> [External Services] -> [Transactions]
-> [Storage: Bulk Operation] -> [Cache: Invalidate/Set] -> [Event Emission]
-> [Search Indexing: Bulk Update] -> [Encode Response] -> [Post-Response Hooks]
-> [Observability] -> [Backup/Recovery] -> [Shutdown] -> [Config Management]
-> [Schema Evolution]
*Bulk Validation: Validate each item but respond with collective errors.
*Transactions: Ensure the entire batch succeeds or rolls back.
*Event Emission: Emit batched or individual events based on use case.
*Search Indexing: Optimize with batch indexing if possible.

*/

// Middleware & Functional Stubs
func TransportSecurity() { fmt.Println("[Transport Security]: Ensuring secure connection") }
func ContextAndTracing() {
	fmt.Println("[Context & Tracing]: Setting up request context and distributed tracing")
}
func RateLimiting()      { fmt.Println("[Rate Limiting]: Checking rate limits") }
func AuthNAuthZ()        { fmt.Println("[AuthN/AuthZ]: Authenticating and authorizing user") }
func IdempotencyChecks() { fmt.Println("[Idempotency Checks]: Ensuring idempotent request handling") }
func Deserialize[T any](input string) ([]T, error) {
	var objs []T
	fmt.Println("[Deserialize]: Deserializing input")
	return objs, nil
}
func Sanitize[T any](obj *T)           { fmt.Println("[Sanitize]: Cleaning input data") }
func Validate[T any](obj *T) error     { fmt.Println("[Validate]: Validating data"); return nil }
func TransformAndEnrich[T any](obj *T) { fmt.Println("[Transform & Enrich]: Enriching data") }
func ExternalServices[T any](obj *T) {
	fmt.Println("[External Services]: Communicating with external dependencies")
}
func Transactions(action func() error) error {
	fmt.Println("[Transactions]: Beginning transaction")
	defer fmt.Println("[Transactions]: Ending transaction")
	return action()
}
func Storage[T any](obj *T, operation string) {
	fmt.Printf("[Storage]: Performing %s operation on persistent storage\n", operation)
}
func Cache[T any](obj *T, operation string) {
	fmt.Printf("[Cache]: Performing %s operation on cache\n", operation)
}
func EventEmission[T any](obj *T) { fmt.Println("[Event Emission]: Emitting event") }
func SearchIndexing[T any](obj *T, operation string) {
	fmt.Printf("[Search Indexing]: Performing %s operation on search index\n", operation)
}
func EncodeResponse[T any](obj *T) { fmt.Println("[Encode Response]: Encoding and sending response") }
func PostResponseHooks()           { fmt.Println("[Post-Response Hooks]: Running post-response actions") }
func Observability()               { fmt.Println("[Observability]: Logging metrics and traces") }
func BackupRecovery()              { fmt.Println("[Backup/Recovery]: Ensuring data safety") }
func Shutdown()                    { fmt.Println("[Shutdown]: Cleaning up resources") }
func ConfigManagement()            { fmt.Println("[Config Management]: Loading and applying configuration") }
func SchemaEvolution()             { fmt.Println("[Schema Evolution]: Managing schema changes") }

// Create supports generics flow from REST to SQL.
// It takes an input string, but should work on a serialized object/reader
func Create[T any](input string) {
	fmt.Printf("\n\n[CREATE]: POST /entity\n")
	TransportSecurity()
	ContextAndTracing()
	RateLimiting()
	AuthNAuthZ()
	IdempotencyChecks()
	// collect errors while supporting bulk operations
	objs, _ := Deserialize[T](input)
	for i := range objs {
		Sanitize(&objs[i])
		_ = Validate(&objs[i])
		TransformAndEnrich(&objs[i])
		ExternalServices(&objs[i])

		_ = Transactions(func() error {
			Storage(&objs[i], "create")
			Cache(&objs[i], "set")
			return nil
		})

		EventEmission(&objs[i])
		SearchIndexing(&objs[i], "add")
		EncodeResponse(&objs[i])
	}

	PostResponseHooks()
	Observability()
	BackupRecovery()
	Shutdown()
	ConfigManagement()
	SchemaEvolution()
}

// List gets requested objects from storage and transports them.
//
// Parameters:
//
//	filter: string.
//
// Returns:
//
//	A list of entries of type T.
func List[T any](filter string) {
	fmt.Printf("\n\n[READ]: GET /entity\n")
	TransportSecurity()
	ContextAndTracing()
	RateLimiting()
	AuthNAuthZ()

	fmt.Println("[List]: Retrieving filtered list of records")
	PostResponseHooks()
	Observability()
	Shutdown()
	ConfigManagement()
	SchemaEvolution()
}

func Read[T any](id string) {
	fmt.Printf("\n\n[READ]: GET /entity/{id}\n")
	TransportSecurity()
	ContextAndTracing()
	RateLimiting()
	AuthNAuthZ()

	var obj T
	Cache(&obj, "read")
	Storage(&obj, "select")
	SearchIndexing(&obj, "query")
	EncodeResponse(&obj)
	PostResponseHooks()
	Observability()
	Shutdown()
	ConfigManagement()
	SchemaEvolution()
}

func Update[T any](input string) {
	fmt.Printf("\n\n[UPDATE]: PUT /entity/{id}\n")
	TransportSecurity()
	ContextAndTracing()
	RateLimiting()
	AuthNAuthZ()
	IdempotencyChecks()
	// collect errors while supporting bulk operations
	objs, _ := Deserialize[T](input)
	for i := range objs {
		Sanitize(&objs[i])
		_ = Validate(&objs[i])
		TransformAndEnrich(&objs[i])
		ExternalServices(&objs[i])

		_ = Transactions(func() error {
			Storage(&objs[i], "update")
			Cache(&objs[i], "update")
			return nil
		})

		EventEmission(&objs[i])
		SearchIndexing(&objs[i], "update")
		EncodeResponse(&objs[i])
	}

	PostResponseHooks()
	Observability()
	BackupRecovery()
	Shutdown()
	ConfigManagement()
	SchemaEvolution()
}

func Delete[T any](id string) {
	fmt.Printf("\n\n[DELETE]: DELETE /entity/{id}\n")
	TransportSecurity()
	ContextAndTracing()
	RateLimiting()
	AuthNAuthZ()
	IdempotencyChecks()

	var obj T
	_ = Transactions(func() error {
		Storage(&obj, "delete")
		Cache(&obj, "invalidate")
		return nil
	})

	EventEmission(&obj)
	SearchIndexing(&obj, "remove")
	EncodeResponse(&obj)
	PostResponseHooks()
	Observability()
	BackupRecovery()
	Shutdown()
	ConfigManagement()
	SchemaEvolution()
}

type HandlerFunc[T any] func(input string)

type RouteRegistry[T any] struct {
	Create HandlerFunc[T]
	Read   func(id string)
	Update HandlerFunc[T]
	Delete func(id string)
	List   func(filter string)
}

func RegisterRoutes[T any](resource string, registry RouteRegistry[T]) {
	fmt.Printf("[Registering Routes] Resource: %s\n", resource)
	fmt.Println("- POST /", resource)
	fmt.Println("- GET /", resource)
	fmt.Println("- GET /", resource, "/{id}")
	fmt.Println("- PUT /", resource, "/{id}")
	fmt.Println("- DELETE /", resource, "/{id}")
}

// QUERY GENERATORS
func GenerateInsertQuery[T any](obj T) string {
	v := reflect.ValueOf(obj)
	typeOfT := v.Type()
	fields := []string{}
	placeholders := []string{}

	for i := 0; i < v.NumField(); i++ {
		fields = append(fields, typeOfT.Field(i).Name)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		typeOfT.Name(),
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)
	return query
}

func GenerateUpdateQuery[T any](obj T) string {
	v := reflect.ValueOf(obj)
	typeOfT := v.Type()
	fields := []string{}
	var idValue string

	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfT.Field(i).Name
		if strings.ToLower(fieldName) == "id" {
			idValue = fmt.Sprintf("$%d", i+1)
			continue
		}
		fields = append(fields, fmt.Sprintf("%s=$%d", fieldName, i+1))
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id=%s",
		typeOfT.Name(),
		strings.Join(fields, ", "),
		idValue,
	)
	return query
}

func GenerateSelectQuery[T any](id string) string {
	var obj T
	typeOfT := reflect.TypeOf(obj)
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", typeOfT.Name())
	return query
}

func GenerateDeleteQuery[T any](id string) string {
	var obj T
	typeOfT := reflect.TypeOf(obj)
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", typeOfT.Name())
	return query
}

// Example entities
type Table struct {
	ID   string
	Name string
}

type Row struct {
	ID      string
	TableID string
	Data    string
}

// Main for example calls
type ExampleEntity struct {
	ID   string
	Name string
}

func main() {
	tableRegistry := RouteRegistry[Table]{
		Create: Create[Table],
		Read:   Read[Table],
		Update: Update[Table],
		Delete: Delete[Table],
		List:   List[Table],
	}

	rowRegistry := RouteRegistry[Row]{
		Create: Create[Row],
		Read:   Read[Row],
		Update: Update[Row],
		Delete: Delete[Row],
		List:   List[Row],
	}

	RegisterRoutes("tables", tableRegistry)
	RegisterRoutes("rows", rowRegistry)

	Create[ExampleEntity]("[{\"Name\": \"Example\"}]")
	Read[ExampleEntity]("123")
	Update[ExampleEntity]("[{\"ID\": \"123\", \"Name\": \"Updated Example\"}]")
	Delete[ExampleEntity]("123")
	List[ExampleEntity]("name=example")

	example := ExampleEntity{ID: "123", Name: "Example"}

	fmt.Println(GenerateInsertQuery(example))
	fmt.Println(GenerateUpdateQuery(example))
	fmt.Println(GenerateSelectQuery[ExampleEntity]("123"))
	fmt.Println(GenerateDeleteQuery[ExampleEntity]("123"))
}
