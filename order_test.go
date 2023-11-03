package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	order := newOrder(Bid, 5)

	assert.Equal(t, Bid, order.side, "order side should be Bid")
	assert.Equal(t, uint64(5), order.size, "order size should be 5")
}
