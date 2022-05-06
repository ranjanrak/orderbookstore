package orderbookstore

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMock(mockRow *sqlmock.Rows, query string) *Client {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectQuery(query).
		WillReturnRows(mockRow)

	cli := &Client{
		dbClient:    db,
		apiKey:      "your_api_key",
		accessToken: "your_access_token",
	}
	return cli
}

func TestQuerySymbol(t *testing.T) {
	// Add mock row for test
	mockedRow := sqlmock.NewRows([]string{"order_timestamp", "order_id", "exchange", "tradingsymbol", "average_price", "transaction_type"}).
		AddRow("2022-05-06 14:27:47", "123456", "NSE", "SBIN", 242.75, "BUY")
	// Add expected query
	query := `SELECT order_timestamp, order_id, exchange, tradingsymbol, average_price, transaction_type FROM orderbook FINAL WHERE 
		(tradingsymbol = 'SBIN' AND status = 'COMPLETE') ORDER BY
		(order_timestamp, order_id)`
	dbMock := setupMock(mockedRow, query)
	symbolBook := dbMock.QuerySymbol("SBIN")
	// Expected output
	expectedSymbolBook := SymbolStore{
		SymbolBook{
			Exchange:        "NSE",
			Symbol:          "SBIN",
			OrderID:         "123456",
			OrderTimestamp:  "2022-05-06 14:27:47",
			AveragePrice:    242.75,
			TransactionType: "BUY",
		},
	}
	assert.Equal(t, expectedSymbolBook, symbolBook, "Actual symbolBook not matching with expectedSymbolBook response")
}
