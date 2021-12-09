package orderbookstore

import (
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go"
)

// Client represents clickhouse DB client connection
type Client struct {
	dbClient *sql.DB
}

// Creates a new DB connection client
func New() *Client {
	// Use DSN as your clickhouse DB setup.
	// visit https://github.com/ClickHouse/clickhouse-go#dsn to know more
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err = connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}
	return &Client{
		dbClient: connect,
	}
}
