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
