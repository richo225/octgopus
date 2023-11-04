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

func TestLimitAddOrder(t *testing.T) {
	limit := newLimit(250)
	order := newOrder(Bid, 5)

	limit.addOrder(order)

	assert.Equal(t, 1, len(limit.orders), "limit should have 1 order")
	assert.Equal(t, order, limit.orders[0], "limit order should be new order")
	assert.Equal(t, uint64(5), limit.totalVolume, "limit total volume should be 5")
}

func TestLimitMatchOrder(t *testing.T) {
	limit := newLimit(250)
	sellOrder := newOrder(Ask, 5)
	buyOrder := newOrder(Bid, 3)

	limit.addOrder(sellOrder)

	matches := limit.matchOrder(buyOrder)

	assert.Equal(t, 1, len(matches), "limit should have 1 match")
	assert.Equal(t, buyOrder, matches[0].bid, "match bid should be order")
	assert.Equal(t, sellOrder, matches[0].ask, "match ask should be order")
	assert.Equal(t, uint64(3), matches[0].sizeFilled, "match size filled should be 3")
	assert.Equal(t, uint64(250), matches[0].price, "match price should be 250")
}

func TestLimitFillOrders(t *testing.T) {
	limit := newLimit(250)
	sellOrder := newOrder(Ask, 1)
	buyOrder := newOrder(Bid, 3)

	limit.addOrder(buyOrder)

	match := limit.fillOrders(buyOrder, sellOrder)

	assert.Equal(t, buyOrder, match.bid, "match bid should be order")
	assert.Equal(t, sellOrder, match.ask, "match ask should be order")
	assert.Equal(t, uint64(1), match.sizeFilled, "match size filled should be 1")
	assert.Equal(t, uint64(250), match.price, "match price should be 250")
}
