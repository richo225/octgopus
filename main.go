package main

import (
	"github.com/kr/pretty"
)

func main() {
	buyOrder := newOrder(Bid, 5)
	sellOrder := newOrder(Ask, 5)
	limit := newLimit(250)
	limit.addOrder(buyOrder)
	limit.addOrder(sellOrder)

	orderbook := newOrderBook()
	orderbook.addOrder(250, buyOrder)
	orderbook.addOrder(200, sellOrder)

	pair := newTradingPair("ETH", "GBP")
	engine := newMatchingEngine()
	engine.addNewMarket(pair)
	engine.orderbooks[pair] = *orderbook

	pretty.Print(engine)
}
