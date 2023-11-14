package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	order := NewOrder(Bid, 5)

	assert.Equal(t, Bid, order.Side, "order side should be Bid")
	assert.Equal(t, float64(5), order.Size, "order size should be 5")
}
