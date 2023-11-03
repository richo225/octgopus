package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderBook(t *testing.T) {
	orderBook := newOrderBook()

	assert.Empty(t, orderBook.asks, "order book should initialize with no asks")
	assert.Empty(t, orderBook.bids, "order book should initialise with no bids")
}

func TestOrderBookPlaceLimitOrder(t *testing.T) {
	orderBook := newOrderBook()
	buyOrder := newOrder(Bid, 5)
	sellOrder := newOrder(Ask, 5)

	orderBook.placeLimitOrder(250, buyOrder)
	orderBook.placeLimitOrder(120, sellOrder)

	assert.Equal(t, 1, len(orderBook.bids), "order book should have 1 limit in bids")
	assert.Equal(t, buyOrder, orderBook.bidLimits[250].orders[0], "order book should have the correct buy order in bidLimits")
	assert.Equal(t, buyOrder, orderBook.bids[0].orders[0], "order book should have the correct buy order in bids")

	assert.Equal(t, 1, len(orderBook.asks), "order book should have 1 limit in asks")
	assert.Equal(t, sellOrder, orderBook.askLimits[120].orders[0], "order book should have the correct sell order in askLimits")
	assert.Equal(t, sellOrder, orderBook.asks[0].orders[0], "order book should have the correct buy order in bids")
}

func TestOrderBookAsks(t *testing.T) {
	orderbook := newOrderBook()
	limit1 := newLimit(12)
	limit2 := newLimit(8)
	limit3 := newLimit(25)

	orderbook.asks = []*Limit{limit1, limit2, limit3}
	sortedAsks := orderbook.Asks()

	assert.Equal(t, 3, len(sortedAsks), "Asks() should return 3 limits")
	assert.Equal(t, []*Limit{limit2, limit1, limit3}, sortedAsks, "Asks() should return limits in ascending price")
}

func TestOrderBookBids(t *testing.T) {
	orderbook := newOrderBook()
	limit1 := newLimit(12)
	limit2 := newLimit(8)
	limit3 := newLimit(25)

	orderbook.bids = []*Limit{limit1, limit2, limit3}
	sortedBids := orderbook.Bids()

	assert.Equal(t, 3, len(sortedBids), "Bids() should return 3 limits")
	assert.Equal(t, []*Limit{limit3, limit1, limit2}, sortedBids, "Bids() should return limits in descending price")
}
