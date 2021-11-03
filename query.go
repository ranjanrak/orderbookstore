package orderbookstore

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go"
)

// represent tradebook fields
type TradeStore struct {
	OrderTimestamp  string
	Exchange        string
	TradingSymbol   string
	AveragePrice    float64
	TransactionType string
}

// SymbolBook represents symbol specific tradebook fields
type SymbolBook struct {
	Symbol          string
	OrderID         string
	OrderTimestamp  string
	AveragePrice    float64
	TransactionType string
}

// Trades is list of trade
type Trades []TradeStore

// SymbolStore is list of SymbolBook trade
type SymbolStore []SymbolBook

func QuerySymbol(tradingSymbol string) SymbolStore {
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

	queryStatement := fmt.Sprintf(`SELECT 
					     order_timestamp, 
					     order_id, 
					     tradingsymbol, 
					     average_price   
					FROM orderbook 
					FINAL 
					WHERE (tradingsymbol = '%s' AND status = 'COMPLETE')
					ORDER BY (order_timestamp, order_id)`, tradingSymbol)

	rows, err := connect.Query(queryStatement)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			orderTimestamp  string
			orderId         string
			symbol          string
			averagePrice    float64
			transactionType string
		)
		if err := rows.Scan(&orderTimestamp, &orderId, &symbol, &averagePrice, &transactionType); err != nil {
			log.Fatal(err)
		}
		symbolTrade := SymbolBook{
			Symbol:          symbol,
			OrderID:         orderId,
			OrderTimestamp:  orderTimestamp,
			AveragePrice:    averagePrice,
			TransactionType: transactionType,
		}
		symbolList = append(symbolList, symbolTrade)
	}
	return symbolList
}

func TradeBook(startTime, endTime time.Time) Trades {
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

	// DB query to retrive orderbook between StartDate and EndDate
	orderBookStatement := fmt.Sprintf(`SELECT 
									order_timestamp,
									exchange,
									tradingsymbol,
									average_price,
									transaction_type
									FROM orderbook
									FINAL
									WHERE (status = 'COMPLETE' AND order_timestamp >= '%s' AND order_timestamp <= '%s')
									ORDER BY (order_timestamp)`, startT, endT)

	rows, err := connect.Query(orderBookStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tradesList Trades
	for rows.Next() {
		var (
			orderTimestamp  string
			exchange        string
			tradingSymbol   string
			averagePrice    float64
			transactionType string
		)
		if err := rows.Scan(&orderTimestamp, &exchange, &tradingSymbol, &averagePrice, &transactionType); err != nil {
			log.Fatal(err)
		}
		trade := TradeStore{
			OrderTimestamp:  orderTimestamp,
			Exchange:        exchange,
			TradingSymbol:   tradingSymbol,
			AveragePrice:    averagePrice,
			TransactionType: transactionType,
		}
		tradesList = append(tradesList, trade)
	}
	return tradesList
}
