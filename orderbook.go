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
	asks      []*Limit
	bids      []*Limit
	askLimits map[uint64]*Limit
	bidLimits map[uint64]*Limit
}

func newOrderBook() *Orderbook {
	return &Orderbook{
		asks:      []*Limit{},
		bids:      []*Limit{},
		askLimits: make(map[uint64]*Limit),
		bidLimits: make(map[uint64]*Limit),
	}
}

func (book *Orderbook) Asks() []*Limit {
	sort.Slice(book.asks, func(i, j int) bool {
		return book.asks[i].price < book.asks[j].price
	})

	return book.asks
}

func (book *Orderbook) Bids() []*Limit {
	sort.Slice(book.bids, func(i, j int) bool {
		return book.bids[i].price > book.bids[j].price
	})

	return book.bids
}

func (book *Orderbook) bestAsk() *Limit {
	if len(book.asks) == 0 {
		return nil
	}

	return book.Asks()[0]
}

func (book *Orderbook) bestBid() *Limit {
	if len(book.bids) == 0 {
		return nil
	}

	return book.Bids()[0]
}

func (book *Orderbook) totalBidVolume() uint64 {
	var total uint64
	for _, limit := range book.bids {
		total += limit.totalVolume
	}

	return total
}

func (book *Orderbook) totalAskVolume() uint64 {
	var total uint64
	for _, limit := range book.asks {
		total += limit.totalVolume
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
			book.bids = append(book.bids, newLimit)
			newLimit.addOrder(order)
		}
	} else {
		limit, ok := book.askLimits[price]
		if ok {
			limit.addOrder(order)
		} else {
			newLimit := newLimit(price)
			book.askLimits[price] = newLimit
			book.asks = append(book.asks, newLimit)
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
		for _, limit := range book.Asks() {
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

		for _, limit := range book.Bids() {
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

func (book *Orderbook) removeLimit(side Side, limit *Limit) {
	if side == Bid {
		// remove the limit from the orderbook bidLimits
		delete(book.bidLimits, limit.price)
		// remove the limit from the orderbook bids
		for i, l := range book.bids {
			if l == limit {
				book.bids = append(book.bids[:i], book.bids[i+1:]...)
				break
			}
		}
		// resort the bids
		sort.Slice(book.bids, func(i, j int) bool {
			return book.bids[i].price > book.bids[j].price
		})
	} else {
		// remove the limit from the orderbook askLimits
		delete(book.askLimits, limit.price)
		// remove the limit from the orderbook asks
		for i, l := range book.asks {
			if l == limit {
				book.asks = append(book.asks[:i], book.asks[i+1:]...)
				break
			}
		}
		// resort the asks
		sort.Slice(book.asks, func(i, j int) bool {
			return book.asks[i].price < book.asks[j].price
		})
	}
}
