package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTradingPair(t *testing.T) {
	tradingPair := newTradingPair("BTC", "USD")

	assert.Equal(t, "BTC", tradingPair.base, "trading pair should have the correct base currency")
	assert.Equal(t, "USD", tradingPair.quote, "trading pair should have the correct quote currency")
}

func TestNewTradingPlatform(t *testing.T) {
	tradingPlatform := newTradingPlatform()

	assert.Empty(t, tradingPlatform.orderbooks, "matching tradingPlatform should initialise with empty orderbooks")
}

func TestTradingPlatformAddNewMarket(t *testing.T) {
	tradingPlatform := newTradingPlatform()
	pair := TradingPair{"BTC", "USD"}

	tradingPlatform.addNewMarket(pair)
	tradingPlatform.addNewMarket(pair)

	assert.Equal(t, 1, len(tradingPlatform.orderbooks), "tradingPlatform should have 1 order book")
	assert.NotNil(t, tradingPlatform.orderbooks[pair], "tradingPlatform should have an order book for pair")
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
