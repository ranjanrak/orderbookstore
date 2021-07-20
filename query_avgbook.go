package orderbookstore

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/ClickHouse/clickhouse-go"
)

// AvgBook represent avg price and qty detail for both Buy and sell for symbol
type AvgBook struct {
	symbol   string
	buy_avg  float64
	buy_qty  float64
	sell_avg float64
	sell_qty float64
}

func QueryAvgPrice(tradingsymbol string, startTime time.Time, endTime time.Time) AvgBook {

	var (
		total_buy        float64
		buy_qty          float64
		total_sell       float64
		sell_qty         float64
		transaction_type string
		symbol           string
		average_price    float64
		filled_quantity  float64
	)

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

	startT := startTime.Format("2006-01-02 15:04:05")

	endT := endTime.Format("2006-01-02 15:04:05")

	query_statement := fmt.Sprintf(`SELECT
										transaction_type, 
										tradingsymbol, 
										average_price,
										filled_quantity   
									FROM orderbook 
									FINAL 
									WHERE (tradingsymbol = '%s' AND order_timestamp >= '%s' AND order_timestamp <= '%s')
									ORDER BY (order_timestamp, transaction_type)`, tradingsymbol, startT, endT)

	rows, err := connect.Query(query_statement)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&transaction_type, &symbol, &average_price, &filled_quantity); err != nil {
			log.Fatal(err)
		}
		// calculate total buy and sell amount and qty
		if transaction_type == "BUY" && average_price != 0 {
			buy_qty = filled_quantity + buy_qty
			total_buy = average_price*filled_quantity + total_buy
		} else if transaction_type == "SELL" && average_price != 0 {
			sell_qty = filled_quantity + sell_qty
			total_sell = average_price*filled_quantity + total_sell
		}
	}
	// calculate buy and sell avg price
	buy_avg := (total_buy / buy_qty)
	sell_avg := (total_sell / sell_qty)

	avgBook := AvgBook{
		symbol:   symbol,
		buy_avg:  math.Round(buy_avg*100) / 100,
		buy_qty:  buy_qty,
		sell_avg: math.Round(sell_avg*100) / 100,
		sell_qty: sell_qty,
	}
	return avgBook
}
