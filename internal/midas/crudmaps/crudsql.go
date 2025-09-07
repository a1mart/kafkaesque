package crudmaps

// validate & sanitize inputs on post, put, patch

var crud2sql = map[string]string{
	"POST":   "INSERT", //or CREATE (for unknown resources)
	"PUT":    "UPDATE", //maybe put updates resource, patch updates instance
	"PATCH":  "UPDATE", //partial
	"GET":    "SELECT",
	"DELETE": "DELETE",
}

const (
	METHOD = iota
	RESOURCE
	INSTANCE
	QUERY_STRING
)
