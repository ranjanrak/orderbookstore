package orderbookstore

import (
	"fmt"
	"log"
	"math"
	"time"
)

// AvgBook represent avg price and qty detail for both Buy and sell for symbol
type AvgBook struct {
	Symbol      string
	BuyAvg      float64
	BuyQty      float64
	SellAvg     float64
	SellQty     float64
	RealizedPnl float64
}

func (c *Client) QueryAvgPrice(tradingSymbol string, startTime time.Time, endTime time.Time) AvgBook {

	var (
		totalBuy        float64
		buyQty          float64
		totalSell       float64
		sellQty         float64
		transactionType string
		symbol          string
		averagePrice    float64
		filledQuantity  float64
	)

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
					ORDER BY (order_timestamp, transaction_type)`, tradingSymbol, startT, endT)

	rows, err := c.dbClient.Query(query_statement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&transactionType, &symbol, &averagePrice, &filledQuantity); err != nil {
			log.Fatal(err)
		}
		// calculate total buy and sell amount and qty
		if transactionType == "BUY" && averagePrice != 0 {
			buyQty = filledQuantity + buyQty
			totalBuy = averagePrice*filledQuantity + totalBuy
		} else if transactionType == "SELL" && averagePrice != 0 {
			sellQty = filledQuantity + sellQty
			totalSell = averagePrice*filledQuantity + totalSell
		}
	}
	// calculate buy and sell avg price
	buyAvg := (totalBuy / buyQty)
	sellAvg := (totalSell / sellQty)

	// calculate realized P&L
	// rough calculation just to show an idea
	realizedPnl := sellQty*sellAvg - sellQty*buyAvg

	avgBook := AvgBook{
		Symbol:      symbol,
		BuyAvg:      math.Round(buyAvg*100) / 100,
		BuyQty:      buyQty,
		SellAvg:     math.Round(sellAvg*100) / 100,
		SellQty:     sellQty,
		RealizedPnl: math.Round(realizedPnl*100) / 100,
	}
	return avgBook
}
