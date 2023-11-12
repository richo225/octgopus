package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLimit(t *testing.T) {
	limit := newLimit(250)

	assert.Equal(t, float64(250), limit.Price, "price should be 250")
	assert.Empty(t, limit.Orders, "limit orders should be empty")
}

func TestLimitAddOrder(t *testing.T) {
	limit := newLimit(250)
	order := newOrder(Bid, 5)

	limit.addOrder(order)

	assert.Equal(t, 1, len(limit.Orders), "limit should have 1 order")
	assert.Equal(t, order, limit.Orders[0], "limit order should be new order")
	assert.Equal(t, float64(5), limit.TotalVolume, "limit total volume should be 5")
}

func TestLimitRemoveOrder(t *testing.T) {
	limit := newLimit(100)

	order1 := newOrder(Ask, 10)
	order2 := newOrder(Ask, 20)
	order3 := newOrder(Ask, 30)

	limit.addOrder(order1)
	limit.addOrder(order2)
	limit.addOrder(order3)

	limit.removeOrder(order2)

	assert.Equal(t, 2, len(limit.Orders), "limit should have 2 orders")
	assert.Equal(t, float64(40), limit.TotalVolume, "limit should have a total volume of 40")
	assert.Equal(t, order1, limit.Orders[0], "limit should have the correct order at index 0")
	assert.Equal(t, order3, limit.Orders[1], "limit should have the correct order at index 1")
	assert.Equal(t, float64(0), order2.Price, "order should have a price of 0")
}

func TestLimitMatchOrder(t *testing.T) {
	limit := newLimit(250)
	sellOrder := newOrder(Ask, 3)
	buyOrder := newOrder(Bid, 3)

	limit.addOrder(sellOrder)

	matches := limit.matchOrder(buyOrder)

	assert.Equal(t, 1, len(matches), "limit should have 1 match")
	assert.Equal(t, buyOrder, matches[0].Bid, "match bid should be order")
	assert.Equal(t, sellOrder, matches[0].Ask, "match ask should be order")
	assert.Equal(t, float64(3), matches[0].SizeFilled, "match size filled should be 3")
	assert.Equal(t, float64(250), matches[0].Price, "match price should be 250")
	assert.Equal(t, float64(0), limit.TotalVolume, "limit should have the correct total volume")

	// test that the sell limit order is removed
	assert.Equal(t, 0, len(limit.Orders), "limit should have 0 orders")

}

func TestLimitFillOrders(t *testing.T) {
	limit := newLimit(250)
	sellOrder := newOrder(Ask, 1)
	buyOrder := newOrder(Bid, 3)

	limit.addOrder(buyOrder)

	match := limit.fillOrders(buyOrder, sellOrder)

	assert.Equal(t, buyOrder, match.Bid, "match bid should be order")
	assert.Equal(t, sellOrder, match.Ask, "match ask should be order")
	assert.Equal(t, float64(1), match.SizeFilled, "match size filled should be 1")
	assert.Equal(t, float64(250), match.Price, "match price should be 250")
}
