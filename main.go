package main

import (
	"github.com/richo225/octgopus/api"
	"github.com/richo225/octgopus/data"
	"github.com/richo225/octgopus/orderbook"
)

func main() {
	p := orderbook.NewTradingPlatform()
	data.SeedOrderBook(p)

	api.Start(p)
}
