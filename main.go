package main

import (
	"github.com/richo225/octgopus/api"
	"github.com/richo225/octgopus/internal/orderbook"
)

func main() {
	p := orderbook.NewTradingPlatform()
	p.SeedData()

	api.Start(p)
}
