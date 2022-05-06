package orderbookstore

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestQueryAvgPrice(t *testing.T) {
	// Add mock row for test
	mockedRow := sqlmock.NewRows([]string{"transaction_type", "tradingsymbol", "average_price", "filled_quantity"}).
		AddRow("BUY", "IOC", 107.41, 39).
		AddRow("SELL", "IOC", 107.61, 39)

	// Add expected query
	query := `SELECT transaction_type, tradingsymbol, average_price, filled_quantity FROM orderbook FINAL WHERE 
	(tradingsymbol = 'IOC' AND order_timestamp >= '2022-05-06 14:04:05' AND order_timestamp <= '2022-05-06 15:04:05')
	ORDER BY (order_timestamp, transaction_type)`

	dbMock := setupMock(mockedRow, query)
	startTime := time.Date(2022, 5, 6, 14, 04, 05, 0, time.UTC)
	endTime := time.Date(2022, 5, 6, 15, 04, 05, 0, time.UTC)
	avgBook := dbMock.QueryAvgPrice("IOC", startTime, endTime)

	// Expected output
	expectedavgBook := AvgBook{
		Symbol:      "IOC",
		BuyAvg:      107.41,
		BuyQty:      39,
		SellAvg:     107.61,
		SellQty:     39,
		RealizedPnl: 7.8,
	}
	assert.Equal(t, expectedavgBook, avgBook, "Actual avgBook not matching with expectedavgBook response")

}
