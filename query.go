package orderbookstore

import (
	"fmt"
	"log"
	"time"
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
	Exchange        string
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

func (c *Client) QuerySymbol(tradingSymbol string) SymbolStore {
	queryStatement := fmt.Sprintf(`SELECT 
					     order_timestamp, 
					     order_id, 
						 exchange,
					     tradingsymbol, 
					     average_price   
					FROM orderbook 
					FINAL 
					WHERE (tradingsymbol = '%s' AND status = 'COMPLETE')
					ORDER BY (order_timestamp, order_id)`, tradingSymbol)

	rows, err := c.dbClient.Query(queryStatement)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var symbolList SymbolStore
	for rows.Next() {
		var (
			orderTimestamp  string
			orderId         string
			exchange        string
			symbol          string
			averagePrice    float64
			transactionType string
		)
		if err := rows.Scan(&orderTimestamp, &orderId, &exchange, &symbol, &averagePrice, &transactionType); err != nil {
			log.Fatal(err)
		}
		symbolTrade := SymbolBook{
			Exchange:        exchange,
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

func (c *Client) TradeBook(startTime, endTime time.Time) Trades {
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

	rows, err := c.dbClient.Query(orderBookStatement)
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
