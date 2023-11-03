package main

import "sort"

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
			newLimit.addOrder(order)
			book.bidLimits[price] = newLimit
			book.bids = append(book.bids, newLimit)
		}
	} else {
		limit, ok := book.askLimits[price]
		if ok {
			limit.addOrder(order)
		} else {
			newLimit := newLimit(price)
			newLimit.addOrder(order)
			book.askLimits[price] = newLimit
			book.asks = append(book.asks, newLimit)
		}
	}
}

func (book *Orderbook) placeMarketOrder(order *Order) {
	// check which side order is
	if order.side == Bid {
		// get all the sorted asks/bids of opposite side
		// iterate through each limit (market order so start with smallest price)
		for _, limit := range book.Asks() {
			// attempt to match the order to the limit
			limit.matchOrder(order)
			// if the order is filled, break
			if order.size == 0 {
				break
			}

		}
	} else {
		for _, limit := range book.Bids() {
			limit.matchOrder(order)
			if order.size == 0 {
				break
			}
		}
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
