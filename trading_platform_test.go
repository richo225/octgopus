package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTradingPair(t *testing.T) {
	tradingPair := newTradingPair("BTC", "USD")

	assert.Equal(t, "BTC", tradingPair.Base, "trading pair should have the correct base currency")
	assert.Equal(t, "USD", tradingPair.Quote, "trading pair should have the correct quote currency")
}

func TestNewTradingPlatform(t *testing.T) {
	tradingPlatform := newTradingPlatform()

	assert.Empty(t, tradingPlatform.Orderbooks, "matching tradingPlatform should initialise with empty orderbooks")
}

func TestTradingPlatformAddNewMarket(t *testing.T) {
	tradingPlatform := newTradingPlatform()
	pair := TradingPair{"BTC", "USD"}

	tradingPlatform.addNewMarket(pair)
	tradingPlatform.addNewMarket(pair)

	assert.Equal(t, 1, len(tradingPlatform.Orderbooks), "tradingPlatform should have 1 order book")
	assert.NotNil(t, tradingPlatform.Orderbooks[pair], "tradingPlatform should have an order book for pair")
	assert.Equal(t, &pair, tradingPlatform.Orderbooks[pair].Market, "orderbook should have the correct market")
}

func TestTradingPlatformGetOrderBook(t *testing.T) {
	tradingPlatform := newTradingPlatform()
	pair := TradingPair{"BTC", "USD"}

	tradingPlatform.addNewMarket(pair)

	orderBook, err := tradingPlatform.getOrderBook(pair)
	assert.NoError(t, err, "getOrderBook should not return an error")
	assert.NotNil(t, orderBook, "getOrderBook should return an order book")

	nonExistantPair := TradingPair{"ETH", "USD"}
	_, err = tradingPlatform.getOrderBook(nonExistantPair)
	assert.Equal(t, &OrderbookNotFoundError{nonExistantPair}, err, "non existent paair should return OrderbookNotFoundError")
}

func TestTrdingPlatformPlaceLimitOrder(t *testing.T) {
	tradingPlatform := newTradingPlatform()
	pair := TradingPair{"BTC", "USD"}

	tradingPlatform.addNewMarket(pair)

	buyOrder1 := newOrder(Bid, 5)
	buyOrder2 := newOrder(Bid, 8)
	buyOrder3 := newOrder(Bid, 13)
	sellOrder := newOrder(Ask, 5)

	tradingPlatform.placeLimitOrder(pair, 250, buyOrder1)
	tradingPlatform.placeLimitOrder(pair, 250, buyOrder2)
	tradingPlatform.placeLimitOrder(pair, 410, buyOrder3)
	tradingPlatform.placeLimitOrder(pair, 120, sellOrder)

	orderbook, _ := tradingPlatform.getOrderBook(pair)

	assert.Equal(t, uint64(26), orderbook.totalBidVolume(), "order book should have the correct total bid volume")

	assert.Equal(t, 2, len(orderbook.Bids), "order book should have 2 limits in bids")
	assert.Equal(t, buyOrder1, orderbook.bidLimits[250].orders[0], "order book should have the correct buy order in bidLimits")
	assert.Equal(t, buyOrder2, orderbook.bidLimits[250].orders[1], "order book should have the correct buy order in bidLimits")
	assert.Equal(t, buyOrder3, orderbook.bidLimits[410].orders[0], "order book should have the correct buy order in bidLimits")
	assert.Equal(t, buyOrder1, orderbook.Bids[0].orders[0], "order book should have the correct buy order in bids")
	assert.Equal(t, buyOrder2, orderbook.Bids[0].orders[1], "order book should have the correct buy order in bids")
	assert.Equal(t, buyOrder3, orderbook.Bids[1].orders[0], "order book should have the correct buy order in bids")

	assert.Equal(t, 1, len(orderbook.Asks), "order book should have 1 limit in asks")
	assert.Equal(t, sellOrder, orderbook.askLimits[120].orders[0], "order book should have the correct sell order in askLimits")
	assert.Equal(t, sellOrder, orderbook.Asks[0].orders[0], "order book should have the correct buy order in bids")
}
