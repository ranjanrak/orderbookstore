# orderbookstore
A package to store daily end-of-the-day complete order book to clickhouse DB. This can be use to fetch historical order book data to calculate p&l related values.

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
    // Create new orderbook instance
    client := orderbookstore.New()

    // Store current orderbook data to clickhouse DB
    client.DataLoad()

    // Fetch all historical order's for the symbol
    symbolBook := client.QuerySymbol("SBIN")
    fmt.Printf("%+v\n", symbolBook)
    
    startTime := time.Date(2021, 6, 29, 9, 41, 0, 0, time.UTC)
    endTime := time.Date(2021, 8, 3, 15, 05, 0, 0, time.UTC)

    // Fetch average buy and sell price for the mentioned symbol and period
    avgBook := client.QueryAvgPrice("IOC", startTime, endTime)
    fmt.Printf("%+v\n", avgBook)

    // Fetch complete tradebook between start and end date
    tradeBook := client.TradeBook(startTime, endTime)
    fmt.Printf("%+v\n", tradeBook)

}
```

## Response
1> Response for `client.QuerySymbol("SBIN")`:
```
[{Exchange:NSE Symbol:SBIN OrderID:210629000202833 OrderTimestamp:2021-06-29T14:45:24+05:30 
AveragePrice:426.65 TransactionType:BUY} {Exchange:NSE Symbol:SBIN OrderID:210629002938256 
OrderTimestamp:2021-06-29T18:50:27+05:30 AveragePrice:424.05 TransactionType:SELL} 
{Exchange:BSE Symbol:SBIN OrderID:210629002940618 OrderTimestamp:2021-06-29T18:50:54+05:30 
AveragePrice:423.95 TransactionType:SELL}]
```
2> Response for `client.QueryAvgPrice("IOC", startTime, endTime)`:
```
{Symbol:IOC BuyAvg:107.41 BuyQty:39 SellAvg:107.61 SellQty:17 RealizedPnl:-22.9}
```

3> Response for `client.TradeBook(startTime, endTime)`:

```
[{OrderTimestamp:2021-06-29T14:45:11+05:30 Exchange:BSE TradingSymbol:IOC 
AveragePrice:111.05 TransactionType:BUY} 
{OrderTimestamp:2021-06-29T14:45:11+05:30 Exchange:NSE TradingSymbol:IOC 
AveragePrice:111.1 TransactionType:SELL}]
```