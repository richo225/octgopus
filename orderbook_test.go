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

func TestOrderbookTotalBidVolume(t *testing.T) {
	orderBook := newOrderBook()

	order1 := newOrder(Bid, 10)
	order2 := newOrder(Bid, 20)
	order3 := newOrder(Bid, 30)

	orderBook.placeLimitOrder(100, order1)
	orderBook.placeLimitOrder(100, order2)
	orderBook.placeLimitOrder(200, order3)

	assert.Equal(t, uint64(60), orderBook.totalBidVolume(), "order book should have the correct total bid volume")
}

func TestOrderbookTotalAskVolume(t *testing.T) {
	orderBook := newOrderBook()

	order1 := newOrder(Ask, 10)
	order2 := newOrder(Ask, 20)
	order3 := newOrder(Ask, 30)

	orderBook.placeLimitOrder(100, order1)
	orderBook.placeLimitOrder(100, order2)
	orderBook.placeLimitOrder(200, order3)

	assert.Equal(t, uint64(60), orderBook.totalAskVolume(), "order book should have the correct total ask volume")
}

func TestOrderBookPlaceLimitOrder(t *testing.T) {
	orderbook := newOrderBook()
	buyOrder1 := newOrder(Bid, 5)
	buyOrder2 := newOrder(Bid, 8)
	buyOrder3 := newOrder(Bid, 13)
	sellOrder := newOrder(Ask, 5)

	orderbook.placeLimitOrder(250, buyOrder1)
	orderbook.placeLimitOrder(250, buyOrder2)
	orderbook.placeLimitOrder(410, buyOrder3)
	orderbook.placeLimitOrder(120, sellOrder)

	assert.Equal(t, uint64(26), orderbook.totalBidVolume(), "order book should have the correct total bid volume")

	assert.Equal(t, 2, len(orderbook.bids), "order book should have 2 limits in bids")
	assert.Equal(t, buyOrder1, orderbook.bidLimits[250].orders[0], "order book should have the correct buy order in bidLimits")
	assert.Equal(t, buyOrder2, orderbook.bidLimits[250].orders[1], "order book should have the correct buy order in bidLimits")
	assert.Equal(t, buyOrder3, orderbook.bidLimits[410].orders[0], "order book should have the correct buy order in bidLimits")
	assert.Equal(t, buyOrder1, orderbook.bids[0].orders[0], "order book should have the correct buy order in bids")

	assert.Equal(t, 1, len(orderbook.asks), "order book should have 1 limit in asks")
	assert.Equal(t, sellOrder, orderbook.askLimits[120].orders[0], "order book should have the correct sell order in askLimits")
	assert.Equal(t, sellOrder, orderbook.asks[0].orders[0], "order book should have the correct buy order in bids")
}

func TestOrderBookPlaceMarketBuyOrder(t *testing.T) {
	orderbook := newOrderBook()
	buyOrder := newOrder(Bid, 3)
	sellOrder := newOrder(Ask, 8)

	orderbook.placeLimitOrder(250, sellOrder)
	expectedMatches := []Match{{
		sellOrder,
		buyOrder,
		3,
		250,
	}}
	actualMatches, _ := orderbook.placeMarketOrder(buyOrder)
	assert.Equal(t, expectedMatches, actualMatches, "placeMarketOrder should return correct matches")

	assert.Equal(t, uint64(5), orderbook.totalAskVolume(), "total ask volume should be 5")
	assert.Equal(t, uint64(0), buyOrder.size, "buy order size should be 0")
	assert.Equal(t, uint64(5), sellOrder.size, "sell order size should be 5")
}

func TestOrderBookPlaceMarketBuyOrderMultiMatch(t *testing.T) {
	orderbook := newOrderBook()
	buyOrder := newOrder(Bid, 3)
	sellOrder1 := newOrder(Ask, 8)
	sellOrder2 := newOrder(Ask, 2)

	orderbook.placeLimitOrder(250, sellOrder1)
	orderbook.placeLimitOrder(250, sellOrder2)
	expectedMatches := []Match{{
		sellOrder1,
		buyOrder,
		3,
		250,
	}}
	actualMatches, _ := orderbook.placeMarketOrder(buyOrder)
	assert.Equal(t, expectedMatches, actualMatches, "placeMarketOrder should return correct matches")

	assert.Equal(t, uint64(0), buyOrder.size, "buy order size should be 0")
	assert.Equal(t, uint64(5), sellOrder1.size, "sell order size should be 5")
	assert.Equal(t, uint64(2), sellOrder2.size, "sell order size should be 2")
	assert.Equal(t, 1, len(orderbook.asks), "order book should still have 1 limit left")
}

func TestOrderBookPlaceMarketBuyOrderMultiPriceLimitMatch(t *testing.T) {
	orderbook := newOrderBook()
	buyOrder := newOrder(Bid, 3)
	sellOrder1 := newOrder(Ask, 8)
	sellOrder2 := newOrder(Ask, 2)

	orderbook.placeLimitOrder(250, sellOrder1)
	orderbook.placeLimitOrder(240, sellOrder2)
	expectedMatches := []Match{
		{
			sellOrder2,
			buyOrder,
			2,
			240,
		},
		{
			sellOrder1,
			buyOrder,
			1,
			250,
		}}
	actualMatches, _ := orderbook.placeMarketOrder(buyOrder)
	assert.Equal(t, expectedMatches, actualMatches, "placeMarketOrder should return correct matches")

	assert.Equal(t, uint64(0), buyOrder.size, "buy order size should be 0")
	assert.Equal(t, uint64(7), sellOrder1.size, "sell order size should be 7")
	assert.Equal(t, uint64(0), sellOrder2.size, "sell order size should be 8")

	assert.Equal(t, 1, len(orderbook.asks), "order book should have 1 limit left")
	assert.Equal(t, uint64(250), orderbook.asks[0].price, "order book should have the correct non-empty limit left")
	assert.Equal(t, 1, len(orderbook.askLimits), "order book should have 1 limit left")
	assert.Equal(t, uint64(250), orderbook.askLimits[250].price, "order book should have the correct non-empty limit left")
}

func TestOrderBookPlaceMarketBuyOrderInsufficientVolume(t *testing.T) {
	orderbook := newOrderBook()
	buyOrder := newOrder(Bid, 3)
	sellOrder := newOrder(Ask, 2)

	orderbook.placeLimitOrder(250, sellOrder)

	_, err := orderbook.placeMarketOrder(buyOrder)
	assert.Equal(t, &InsufficientVolumeError{2, 3}, err, "placeMarketOrder should return InsufficientVolumeError")
	assert.Equal(t, uint64(3), buyOrder.size, "buy order size should be 3")
}

// ffffffffff

func TestOrderBookPlaceMarketSellOrder(t *testing.T) {
	orderbook := newOrderBook()
	sellOrder := newOrder(Ask, 3)
	buyOrder := newOrder(Bid, 8)

	orderbook.placeLimitOrder(250, buyOrder)
	expectedMatches := []Match{{
		sellOrder,
		buyOrder,
		3,
		250,
	}}
	actualMatches, _ := orderbook.placeMarketOrder(sellOrder)
	assert.Equal(t, expectedMatches, actualMatches, "placeMarketOrder should return correct matches")

	assert.Equal(t, uint64(5), buyOrder.size, "buy order size should be 5")
	assert.Equal(t, uint64(0), sellOrder.size, "sell order size should be 0")
}

func TestOrderBookPlaceMarketSellOrderMultiMatch(t *testing.T) {
	orderbook := newOrderBook()
	sellOrder := newOrder(Ask, 3)
	buyOrder1 := newOrder(Bid, 8)
	buyOrder2 := newOrder(Bid, 2)

	orderbook.placeLimitOrder(250, buyOrder1)
	orderbook.placeLimitOrder(250, buyOrder2)
	expectedMatches := []Match{{
		sellOrder,
		buyOrder1,
		3,
		250,
	}}
	actualMatches, _ := orderbook.placeMarketOrder(sellOrder)
	assert.Equal(t, expectedMatches, actualMatches, "placeMarketOrder should return correct matches")

	assert.Equal(t, uint64(0), sellOrder.size, "sell order size should be 0")
	assert.Equal(t, uint64(5), buyOrder1.size, "buy order size should be 5")
	assert.Equal(t, uint64(2), buyOrder2.size, "sell order size should be 2")
	assert.Equal(t, 1, len(orderbook.bids), "order book should still have 1 limit left")
}

func TestOrderBookPlaceMarketSellOrderMultiPriceLimitMatch(t *testing.T) {
	orderbook := newOrderBook()
	sellOrder := newOrder(Ask, 3)
	buyOrder1 := newOrder(Bid, 8)
	buyOrder2 := newOrder(Bid, 2)

	orderbook.placeLimitOrder(250, buyOrder2)
	orderbook.placeLimitOrder(240, buyOrder1)
	expectedMatches := []Match{
		{
			sellOrder,
			buyOrder2,
			2,
			250,
		},
		{
			sellOrder,
			buyOrder1,
			1,
			240,
		}}
	actualMatches, _ := orderbook.placeMarketOrder(sellOrder)
	assert.Equal(t, expectedMatches, actualMatches, "placeMarketOrder should return correct matches")

	assert.Equal(t, uint64(0), sellOrder.size, "sell order size should be 0")
	assert.Equal(t, uint64(7), buyOrder1.size, "buy order size should be 7")
	assert.Equal(t, uint64(0), buyOrder2.size, "buy order size should be 0")

	assert.Equal(t, 1, len(orderbook.bids), "order book should have 1 limit left")
	assert.Equal(t, uint64(240), orderbook.bids[0].price, "order book should have the correct non-empty limit left")
	assert.Equal(t, 1, len(orderbook.bidLimits), "order book should have 1 limit left")
	assert.Equal(t, uint64(240), orderbook.bidLimits[240].price, "order book should have the correct non-empty limit left")
}

func TestOrderBookPlaceMarketSellOrderInsufficientVolume(t *testing.T) {
	orderbook := newOrderBook()
	sellOrder := newOrder(Ask, 3)
	buyOrder := newOrder(Bid, 2)

	orderbook.placeLimitOrder(250, buyOrder)

	_, err := orderbook.placeMarketOrder(sellOrder)
	assert.Equal(t, &InsufficientVolumeError{2, 3}, err, "placeMarketOrder should return InsufficientVolumeError")
	assert.Equal(t, uint64(3), sellOrder.size, "sell order size should be 3")
}

func TestOrderbookCancelOrder(t *testing.T) {
	orderBook := newOrderBook()

	order1 := newOrder(Bid, 10)
	order2 := newOrder(Bid, 20)
	order3 := newOrder(Bid, 30)

	orderBook.placeLimitOrder(100, order1)
	orderBook.placeLimitOrder(100, order2)
	orderBook.placeLimitOrder(200, order3)

	orderBook.cancelOrder(order2)

	assert.Equal(t, 2, len(orderBook.bids), "order book should have 2 limits in bids")
	assert.Equal(t, uint64(40), orderBook.totalBidVolume(), "order book should have the correct total bid volume")
	assert.Equal(t, 1, len(orderBook.bidLimits[100].orders), "limit should have 1 order")
	assert.Equal(t, uint64(10), orderBook.bidLimits[100].orders[0].size, "limit should have the correct size for order1")
}
