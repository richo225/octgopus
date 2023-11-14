package main

import (
	"github.com/richo225/octgopus/api"
	"github.com/richo225/octgopus/orderbook"
)

func main() {
	p := orderbook.NewTradingPlatform()
	seedData(p)

	api.Start(p)
}
