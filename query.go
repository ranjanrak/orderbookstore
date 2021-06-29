package orderbookstore

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go"
)

func QueryDB(tradingsymbol string) {
	// Use DSN as your clickhouse DB setup.
	// visit https://github.com/ClickHouse/clickhouse-go#dsn to know more
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}
	// example query
	query_statement := fmt.Sprintf(`SELECT 
					     order_timestamp, 
					     order_id, 
					     tradingsymbol, 
					     average_price   
					FROM orderbook 
					FINAL 
					WHERE (tradingsymbol = '%s')
					ORDER BY (order_timestamp, order_id)`, tradingsymbol)

	rows, err := connect.Query(query_statement)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			order_timestamp string
			order_id        string
			symbol          string
			average_price   float64
		)
		if err := rows.Scan(&order_timestamp, &order_id, &symbol, &average_price); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("order_timestamp: %s, order_id: %s, tradingsymbol: %s, average_price: %f\n",
			order_timestamp, order_id, symbol, average_price)
	}
}
