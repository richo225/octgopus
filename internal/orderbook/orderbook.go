package orderbook

import (
	"sort"
)

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

func (book *Orderbook) GetAsks() []*Limit {
	sort.Slice(book.Asks, func(i, j int) bool {
		return book.Asks[i].Price < book.Asks[j].Price
	})

	return book.Asks
}

func (book *Orderbook) GetBids() []*Limit {
	sort.Slice(book.Bids, func(i, j int) bool {
		return book.Bids[i].Price > book.Bids[j].Price
	})

	return book.Bids
}

func (book *Orderbook) bestAsk() *Limit {
	if len(book.Asks) == 0 {
		return nil
	}

	return book.GetAsks()[0]
}

func (book *Orderbook) bestBid() *Limit {
	if len(book.Bids) == 0 {
		return nil
	}

	return book.GetBids()[0]
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
		limit, ok := book.bidLimits[price]
		if ok {
			limit.addOrder(order)
		} else {
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

	if order.Side == Bid {
		if book.totalAskVolume() < order.Size {
			return nil, &InsufficientVolumeError{book.totalAskVolume(), order.Size}
		}

		for _, limit := range book.GetAsks() {
			limitMatches := limit.matchOrder(order)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				book.removeLimit(Ask, limit)
			}

			if order.Size == 0 {
				break
			}

		}
	} else {
		if book.totalBidVolume() < order.Size {
			return nil, &InsufficientVolumeError{book.totalBidVolume(), order.Size}
		}

		for _, limit := range book.GetBids() {
			limitMatches := limit.matchOrder(order)
			matches = append(matches, limitMatches...)

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
		delete(book.bidLimits, limit.Price)

		for i, l := range book.Bids {
			if l == limit {
				book.Bids = append(book.Bids[:i], book.Bids[i+1:]...)
				break
			}
		}
		sort.Slice(book.Bids, func(i, j int) bool {
			return book.Bids[i].Price > book.Bids[j].Price
		})
	} else {
		delete(book.askLimits, limit.Price)

		for i, l := range book.Asks {
			if l == limit {
				book.Asks = append(book.Asks[:i], book.Asks[i+1:]...)
				break
			}
		}
		sort.Slice(book.Asks, func(i, j int) bool {
			return book.Asks[i].Price < book.Asks[j].Price
		})
	}
}
