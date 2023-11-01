package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLimit(t *testing.T) {
	limit := newLimit(250)

	assert.Equal(t, uint64(250), limit.price, "price should be 250")
	assert.Empty(t, limit.orders, "limit orders should be empty")
}

func TestNewOrder(t *testing.T) {
	order := newOrder(Bid, 5)

	assert.Equal(t, Bid, order.side, "order side should be Bid")
	assert.Equal(t, uint64(5), order.size, "order size should be 5")
}
func TestNewOrderBook(t *testing.T) {
	orderBook := newOrderBook()

	assert.Empty(t, orderBook.asks, "order book should initialize with no asks")
	assert.Empty(t, orderBook.bids, "order book should initialise with no bids")
}

func TestLimitAddOrder(t *testing.T) {
	limit := newLimit(250)
	order := newOrder(Bid, 5)

	limit.addOrder(*order)

	assert.Equal(t, 1, len(limit.orders), "limit should have 1 order")
	assert.Equal(t, uint64(5), limit.orders[0].size, "limit order size should be 5")
}

func TestOrderBookAddOrder(t *testing.T) {
	orderBook := newOrderBook()
	buyOrder := newOrder(Bid, 5)
	sellOrder := newOrder(Ask, 5)

	orderBook.addOrder(250, *buyOrder)
	orderBook.addOrder(120, *sellOrder)

	assert.Equal(t, 1, len(orderBook.bids), "order book should have 1 limit in bids")
	assert.Equal(t, *buyOrder, orderBook.bids[250].orders[0], "order book should have the correct buy order in bids")
	assert.Equal(t, 1, len(orderBook.asks), "order book should have 1 limit in asks")
	assert.Equal(t, *sellOrder, orderBook.asks[120].orders[0], "order book should have the correct sell order in asks")
}
