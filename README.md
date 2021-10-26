# orderbookstore
A package to store daily end-of-the-day complete order book to clickhouse DB. This can be used to fetch historical order book data in future to calculate other p&l related data.

## Installation
```
go get -u github.com/ranjanrak/orderbookstore
```

## Usage
```go
package main

import (
    "github.com/ranjanrak/orderbookstore"
)

func main() {
    // Store current orderbook data to clickhouse DB
    orderbookstore.DataLoad()

    // Fetch all historical order's for the requested symbol
    orderbookstore.QueryDB("SBIN")
    
    startTime := time.Date(2021, 6, 29, 9, 41, 0, 0, time.UTC)
    endTime := time.Date(2021, 8, 3, 15, 05, 0, 0, time.UTC)

    // Fetch average buy and sell price for the mentioned symbol and period
    avgBook := orderbookstore.QueryAvgPrice("IOC", startTime, endTime)
    fmt.Printf("%+v\n", avgBook)

    // Fetch complete tradebook between start and end date
    tradeBook := orderbookstore.TradeBook(startTime, endTime)
    fmt.Printf("%+v\n", tradeBook)

}
```

## Response
1> Response for `orderbookstore.QueryDB("SBIN")`:
```
orderTimestamp: 2021-06-23 09:15:14 +0000 UTC, orderId: XXXXXXX, tradingSymbol: SBIN, averagePrice: 421.800000
orderTimestamp: 2021-06-23 09:16:48 +0000 UTC, orderId: XXXXXXX, tradingSymbol: SBIN, averagePrice: 421.500000
orderTimestamp: 2021-06-25 09:07:46 +0000 UTC, orderId: XXXXXXX, tradingSymbol: SBIN, averagePrice: 421.000000
orderTimestamp: 2021-06-25 09:15:08 +0000 UTC, orderId: XXXXXXX, tradingSymbol: SBIN, averagePrice: 420.250000
```
2> Response for `orderbookstore.QueryAvgPrice("IOC", startTime, endTime)`:
```
{Symbol:IOC BuyAvg:107.41 BuyQty:39 SellAvg:107.61 SellQty:17 RealizedPnl:-22.9}
```

3> Response for `orderbookstore.TradeBook(startTime, endTime)`:

```
[{OrderTimestamp:2021-06-29T14:45:11+05:30 Exchange:BSE TradingSymbol:IOC 
AveragePrice:111.05 TransactionType:BUY} 
{OrderTimestamp:2021-06-29T14:45:11+05:30 Exchange:NSE TradingSymbol:IOC 
AveragePrice:111.1 TransactionType:SELL}]
```