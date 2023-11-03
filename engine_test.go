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
