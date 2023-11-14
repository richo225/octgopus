package orderbook

import "fmt"

type InsufficientVolumeError struct {
	available float64
	requested float64
}

func (e *InsufficientVolumeError) Error() string {
	return "InsufficientVolume : " + fmt.Sprint(e.available) + " < " + fmt.Sprint(e.requested)
}

type OrderbookNotFoundError struct {
	pair TradingPair
}

func (e *OrderbookNotFoundError) Error() string {
	return "MarketNotfound : " + e.pair.ToString()
}
