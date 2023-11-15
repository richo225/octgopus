package orderbook

import (
	"fmt"
	"net/http"
)

type InsufficientVolumeError struct {
	available float64
	requested float64
}

func (e *InsufficientVolumeError) Error() string {
	return "InsufficientVolume : " + fmt.Sprint(e.available) + " < " + fmt.Sprint(e.requested)
}

func (e *InsufficientVolumeError) HTTPCode() int {
	return http.StatusBadRequest
}

type OrderbookNotFoundError struct {
	pair TradingPair
}

func (e *OrderbookNotFoundError) Error() string {
	return "MarketNotfound : " + e.pair.ToString()
}

func (e *OrderbookNotFoundError) HTTPCode() int {
	return http.StatusNotFound
}
