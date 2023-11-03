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

func TestOrderBookAddOrder(t *testing.T) {
	orderBook := newOrderBook()
	buyOrder := newOrder(Bid, 5)
	sellOrder := newOrder(Ask, 5)

	orderBook.addOrder(250, buyOrder)
	orderBook.addOrder(120, sellOrder)

	assert.Equal(t, 1, len(orderBook.bids), "order book should have 1 limit in bids")
	assert.Equal(t, buyOrder, orderBook.bids[250].orders[0], "order book should have the correct buy order in bids")
	assert.Equal(t, 1, len(orderBook.asks), "order book should have 1 limit in asks")
	assert.Equal(t, sellOrder, orderBook.asks[120].orders[0], "order book should have the correct sell order in asks")
}
