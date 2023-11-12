package main

import (
	"fmt"
	"sort"
)

type InsufficientVolumeError struct {
	available float64
	requested float64
}

func (e *InsufficientVolumeError) Error() string {
	return "InsufficientVolume : " + fmt.Sprint(e.available) + " < " + fmt.Sprint(e.requested)
}

type Orderbook struct {
	Market    *TradingPair       `json:"market"`
	Asks      []*Limit           `json:"asks"`
	Bids      []*Limit           `json:"bids"`
	askLimits map[float64]*Limit `json:"-"`
	bidLimits map[float64]*Limit `json:"-"`
}

func newOrderBook() *Orderbook {
	return &Orderbook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		askLimits: make(map[float64]*Limit),
		bidLimits: make(map[float64]*Limit),
	}
}

func (book *Orderbook) getAsks() []*Limit {
	sort.Slice(book.Asks, func(i, j int) bool {
		return book.Asks[i].Price < book.Asks[j].Price
	})

	return book.Asks
}

func (book *Orderbook) getBids() []*Limit {
	sort.Slice(book.Bids, func(i, j int) bool {
		return book.Bids[i].Price > book.Bids[j].Price
	})

	return book.Bids
}

func (book *Orderbook) bestAsk() *Limit {
	if len(book.Asks) == 0 {
		return nil
	}

	return book.getAsks()[0]
}

func (book *Orderbook) bestBid() *Limit {
	if len(book.Bids) == 0 {
		return nil
	}

	return book.getBids()[0]
}

func (book *Orderbook) totalBidVolume() float64 {
	var total float64
	for _, limit := range book.Bids {
		total += limit.TotalVolume
	}

	return total
}

func (book *Orderbook) totalAskVolume() float64 {
	var total float64
	for _, limit := range book.Asks {
		total += limit.TotalVolume
	}

	return total
}

func (book *Orderbook) placeLimitOrder(price float64, order *Order) *Order {
	if order.Side == Bid {
		// check if a limit already exists at the order price
		limit, ok := book.bidLimits[price]
		if ok {
			// if one exists, add the order to that limit
			limit.addOrder(order)
		} else {
			// if doesn't exist, create a new limit with the order
			newLimit := newLimit(price)
			book.bidLimits[price] = newLimit
			book.Bids = append(book.Bids, newLimit)
			newLimit.addOrder(order)
		}
	} else {
		limit, ok := book.askLimits[price]
		if ok {
			limit.addOrder(order)
		} else {
			newLimit := newLimit(price)
			newLimit.addOrder(order)
			book.askLimits[price] = newLimit
			book.Asks = append(book.Asks, newLimit)
		}
	}

	return order
}

func (book *Orderbook) placeMarketOrder(order *Order) ([]Match, error) {
	matches := []Match{}

	// check which side order is
	if order.Side == Bid {
		// if total ask volume is less than order size, return error
		if book.totalAskVolume() < order.Size {
			return nil, &InsufficientVolumeError{book.totalAskVolume(), order.Size}
		}

		// get all the sorted asks/bids of opposite side
		// iterate through each limit (market order so start with smallest price)
		for _, limit := range book.getAsks() {
			// attempt to match the order to the limit
			limitMatches := limit.matchOrder(order)
			matches = append(matches, limitMatches...)

			// if the limit is empty, remove it from the orderbook
			if len(limit.Orders) == 0 {
				book.removeLimit(Ask, limit)
			}

			// if the order is filled, break
			if order.Size == 0 {
				break
			}

		}
	} else {
		// if total bid volume is less than order size, return error
		if book.totalBidVolume() < order.Size {
			return nil, &InsufficientVolumeError{book.totalBidVolume(), order.Size}
		}

		for _, limit := range book.getBids() {
			limitMatches := limit.matchOrder(order)
			matches = append(matches, limitMatches...)

			// if the limit is empty, remove it from the orderbook
			if len(limit.Orders) == 0 {
				book.removeLimit(Bid, limit)
			}

			if order.Size == 0 {
				break
			}
		}
	}

	return matches, nil
}

func (book *Orderbook) cancelOrder(order *Order) {
	var limit *Limit

	price := order.Price
	side := order.Side
	if side == Bid {
		limit = book.bidLimits[price]
	} else {
		limit = book.askLimits[price]
	}
	if limit == nil {
		return
	}

	limit.removeOrder(order)

	if len(limit.Orders) == 0 {
		book.removeLimit(side, limit)
	}
}

func (book *Orderbook) removeLimit(side Side, limit *Limit) {
	if side == Bid {
		// remove the limit from the orderbook bidLimits
		delete(book.bidLimits, limit.Price)
		// remove the limit from the orderbook bids
		for i, l := range book.Bids {
			if l == limit {
				book.Bids = append(book.Bids[:i], book.Bids[i+1:]...)
				break
			}
		}
		// resort the bids
		sort.Slice(book.Bids, func(i, j int) bool {
			return book.Bids[i].Price > book.Bids[j].Price
		})
	} else {
		// remove the limit from the orderbook askLimits
		delete(book.askLimits, limit.Price)
		// remove the limit from the orderbook asks
		for i, l := range book.Asks {
			if l == limit {
				book.Asks = append(book.Asks[:i], book.Asks[i+1:]...)
				break
			}
		}
		// resort the asks
		sort.Slice(book.Asks, func(i, j int) bool {
			return book.Asks[i].Price < book.Asks[j].Price
		})
	}
}
