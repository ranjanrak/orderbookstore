# orderbookstore
A utility to store daily end-of-the-day complete order book to clickhouse DB. This can be used to fetch historical order book data in future.

## Installation
```
go get -u github.com/ranjanrak/orderbookstore
```

## Usage
```
package main

import (
    "fmt"
    "github.com/ranjanrak/orderbookstore"
)

func main() {
    data := orderbookstore.DataLoad()
	fmt.Printf("%+v\n", data)
}
```