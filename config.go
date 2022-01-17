package orderbookstore

import (
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go"
)

// Client represents clickhouse DB client connection
type Client struct {
	dbClient    *sql.DB
	apiKey      string
	accessToken string
}

// ClientParam represents interface to connect clickhouse and kiteconnect API
type ClientParam struct {
	DBSource    string
	ApiKey      string
	AccessToken string
}

// Creates a new DB connection client
func New(userParam ClientParam) *Client {
	if userParam.DBSource == "" {
		userParam.DBSource = "tcp://127.0.0.1:9000?debug=true"
	}
	connect, err := sql.Open("clickhouse", userParam.DBSource)
	if err = connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}
	return &Client{
		dbClient:    connect,
		apiKey:      userParam.ApiKey,
		accessToken: userParam.AccessToken,
	}
}
