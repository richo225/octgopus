package main

import (
	"fmt"
	"sort"
)

type InsufficientVolumeError struct {
	available uint64
	requested uint64
}

func (e *InsufficientVolumeError) Error() string {
	return "InsufficientVolume : " + fmt.Sprint(e.available) + " < " + fmt.Sprint(e.requested)
}

type Orderbook struct {
	Market    *TradingPair      `json:"market"`
	Asks      []*Limit          `json:"asks"`
	Bids      []*Limit          `json:"bids"`
	askLimits map[uint64]*Limit `json:"-"`
	bidLimits map[uint64]*Limit `json:"-"`
}

func newOrderBook() *Orderbook {
	return &Orderbook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		askLimits: make(map[uint64]*Limit),
		bidLimits: make(map[uint64]*Limit),
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

func (book *Orderbook) totalBidVolume() uint64 {
	var total uint64
	for _, limit := range book.Bids {
		total += limit.TotalVolume
	}

	return total
}

func (book *Orderbook) totalAskVolume() uint64 {
	var total uint64
	for _, limit := range book.Asks {
		total += limit.TotalVolume
	}

	return total
}

func (book *Orderbook) placeLimitOrder(price uint64, order *Order) {
	if order.side == Bid {
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
			book.askLimits[price] = newLimit
			book.Asks = append(book.Asks, newLimit)
			newLimit.addOrder(order)
		}
	}
}

func (book *Orderbook) placeMarketOrder(order *Order) ([]Match, error) {
	matches := []Match{}

	// check which side order is
	if order.side == Bid {
		// if total ask volume is less than order size, return error
		if book.totalAskVolume() < order.size {
			return nil, &InsufficientVolumeError{book.totalAskVolume(), order.size}
		}

		// get all the sorted asks/bids of opposite side
		// iterate through each limit (market order so start with smallest price)
		for _, limit := range book.getAsks() {
			// attempt to match the order to the limit
			limitMatches := limit.matchOrder(order)
			matches = append(matches, limitMatches...)

			// if the limit is empty, remove it from the orderbook
			if len(limit.orders) == 0 {
				book.removeLimit(Ask, limit)
			}

			// if the order is filled, break
			if order.size == 0 {
				break
			}

		}
	} else {
		// if total bid volume is less than order size, return error
		if book.totalBidVolume() < order.size {
			return nil, &InsufficientVolumeError{book.totalBidVolume(), order.size}
		}

		for _, limit := range book.getBids() {
			limitMatches := limit.matchOrder(order)
			matches = append(matches, limitMatches...)

			// if the limit is empty, remove it from the orderbook
			if len(limit.orders) == 0 {
				book.removeLimit(Bid, limit)
			}

			if order.size == 0 {
				break
			}
		}
	}

	return matches, nil
}

func (book *Orderbook) cancelOrder(order *Order) {
	limit := order.limit
	side := order.side
	if limit == nil {
		return
	}

	limit.removeOrder(order)

	if len(limit.orders) == 0 {
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
