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

func TestNewMatchingEngine(t *testing.T) {
	matchingEngine := newMatchingEngine()

	assert.Empty(t, matchingEngine.orderbooks, "matching engine should initialise with empty orderbooks")
}

func TestMatchingEngineAddNewMarket(t *testing.T) {
	engine := newMatchingEngine()
	pair := TradingPair{"BTC", "USD"}

	engine.addNewMarket(pair)
	engine.addNewMarket(pair)

	assert.Equal(t, 1, len(engine.orderbooks), "engine should have 1 order book")
	assert.NotNil(t, engine.orderbooks[pair], "engine should have an order book for pair")
}

func TestMatchingEngineGetOrderBook(t *testing.T) {
	engine := newMatchingEngine()
	pair := TradingPair{"BTC", "USD"}

	engine.addNewMarket(pair)

	orderBook, err := engine.getOrderBook(pair)
	assert.NoError(t, err, "getOrderBook should not return an error")
	assert.NotNil(t, orderBook, "getOrderBook should return an order book")

	nonExistantPair := TradingPair{"ETH", "USD"}
	_, err = engine.getOrderBook(nonExistantPair)
	assert.Equal(t, &OrderbookNotFoundError{nonExistantPair}, err, "non existent paair should return OrderbookNotFoundError")
}
