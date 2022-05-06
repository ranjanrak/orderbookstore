package orderbookstore

import (
	"log"

	kiteconnect "github.com/zerodha/gokiteconnect/v3"
)

func (c *Client) DataLoad() {
	// Create a new Kite connect instance
	kc := kiteconnect.New(c.apiKey)

	// Set access token
	kc.SetAccessToken(c.accessToken)

	_, err := c.dbClient.Exec(`
		CREATE TABLE IF NOT EXISTS orderbook (
			order_id           LowCardinality(String),
			parent_order_id    LowCardinality(String),
			exchange_order_id  LowCardinality(String),
			placed_by          LowCardinality(String),
			variety            LowCardinality(String),
			status             LowCardinality(String),
			tradingsymbol      LowCardinality(String),
			exchange           LowCardinality(String),
			instrument_token   UInt32,
			transaction_type   LowCardinality(String),
			order_type         LowCardinality(String),
			product            LowCardinality(String),
			validity           LowCardinality(String),
			price              FLOAT(),
			quantity           FLOAT(),
			trigger_price      FLOAT(),
			average_price      FLOAT(),
			pending_quantity   FLOAT(),
			filled_quantity    FLOAT(),
			disclosed_quantity FLOAT(),
			order_timestamp    DateTime('Asia/Calcutta'),
			exchange_timestamp DateTime('Asia/Calcutta'),
			status_message     LowCardinality(String),
			tag                LowCardinality(String)
		) engine=ReplacingMergeTree()
		PARTITION BY toYYYYMMDD(order_timestamp)
		ORDER BY (order_id, order_timestamp, tradingsymbol)
	`)

	if err != nil {
		log.Fatal(err)
	}

	tx, err := c.dbClient.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`INSERT INTO orderbook (order_id, parent_order_id, exchange_order_id,
				placed_by, variety, status, tradingsymbol, exchange, instrument_token, transaction_type,
				order_type, product, validity, price, quantity, trigger_price, average_price, pending_quantity,
				filled_quantity, disclosed_quantity, order_timestamp, exchange_timestamp,
				status_message, tag) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Fatal(err)
	}

	orders, err := kc.GetOrders()
	if err != nil {
		log.Fatal(err)
	}

	for _, order := range orders {
		if _, err := stmt.Exec(
			order.OrderID,
			order.ParentOrderID,
			order.ExchangeOrderID,
			order.PlacedBy,
			order.Variety,
			order.Status,
			order.TradingSymbol,
			order.Exchange,
			order.InstrumentToken,
			order.TransactionType,
			order.OrderType,
			order.Product,
			order.Validity,
			order.Price,
			order.Quantity,
			order.TriggerPrice,
			order.AveragePrice,
			order.PendingQuantity,
			order.FilledQuantity,
			order.DisclosedQuantity,
			order.OrderTimestamp.Time,
			order.ExchangeTimestamp.Time,
			order.StatusMessage,
			order.Tag,
		); err != nil {
			log.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
