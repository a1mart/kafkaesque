// connectors/sinks/postgres/postgres_sink.go
package postgres

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    connectors "github.com/a1mart/kafkaesque/internal/indranet"
)

type PostgresSink struct {
    db *sql.DB
}

func (p *PostgresSink) Init(config map[string]string) error {
    connStr, ok := config["dsn"]
    if !ok {
        return fmt.Errorf("missing 'dsn' config")
    }

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return err
    }

    p.db = db
    fmt.Println("PostgreSQL sink initialized")
    return nil
}

func (p *PostgresSink) Start() error {
    fmt.Println("PostgreSQL sink started")
    return nil
}

func (p *PostgresSink) Stop() error {
    fmt.Println("Closing PostgreSQL connection")
    return p.db.Close()
}

// Factory function
func NewPostgresSink() connectors.Connector {
    return &PostgresSink{}
}
