// connectors/api/handlers.go
package api

import (
	"encoding/json"
	"net/http"

	connectors "github.com/a1mart/kafkaesque/internal/indranet"
	postgres "github.com/a1mart/kafkaesque/internal/indranet/sinks"
	kafka "github.com/a1mart/kafkaesque/internal/indranet/sources"
)

var manager = connectors.NewConnectorManager()

type RegisterRequest struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

// Registers a new connector
func RegisterConnector(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var connector connectors.Connector
	switch req.Type {
	case "kafka_source":
		connector = kafka.NewKafkaSource()
	case "postgres_sink":
		connector = postgres.NewPostgresSink()
	default:
		http.Error(w, "Unknown connector type", http.StatusBadRequest)
		return
	}

	if err := manager.Register(req.Name, connector, req.Config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Starts a connector
func StartConnector(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if err := manager.Start(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Stops a connector
func StopConnector(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if err := manager.Stop(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Lists all registered connectors
func ListConnectors(w http.ResponseWriter, r *http.Request) {
	connectors := manager.ListConnectors()
	json.NewEncoder(w).Encode(connectors)
}

/*
// main.go
package main

import (
    "github.com/a1mart/kafkaesque/connectors/api"
    "net/http"
)

func main() {
    http.HandleFunc("/connectors/register", api.RegisterConnector)
    http.HandleFunc("/connectors/start", api.StartConnector)
    http.HandleFunc("/connectors/stop", api.StopConnector)
    http.HandleFunc("/connectors/list", api.ListConnectors)

    http.ListenAndServe(":8080", nil)
}

curl -X POST http://localhost:8080/connectors/register -d '{
    "name": "kafka-source-1",
    "type": "kafka_source",
    "config": { "brokers": "localhost:9092", "topic": "events" }
}' -H "Content-Type: application/json"
curl -X GET "http://localhost:8080/connectors/start?name=kafka-source-1"
curl -X GET "http://localhost:8080/connectors/list"
*/
