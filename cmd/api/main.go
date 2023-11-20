package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/richo225/octgopus/internal/api"
	"github.com/richo225/octgopus/internal/orderbook"
)

func main() {
	p := orderbook.NewTradingPlatform()
	p.SeedData()

	api.Start(p)
}
