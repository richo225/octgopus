package main

import "sort"

type Match struct {
	ask        *Order
	bid        *Order
	sizeFilled uint64
	price      uint64
}

type Limit struct {
	Price       uint64   `json:"price"`
	TotalVolume uint64   `json:"total_volume"`
	orders      []*Order `json:"-"`
}

func newLimit(price uint64) *Limit {
	return &Limit{
		Price:  price,
		orders: []*Order{},
	}
}

func (limit *Limit) addOrder(order *Order) {
	order.limit = limit
	limit.orders = append(limit.orders, order)
	limit.TotalVolume += order.size
}

func (limit *Limit) removeOrder(order *Order) {
	for i, o := range limit.orders {
		if o == order {
			limit.orders = append(limit.orders[:i], limit.orders[i+1:]...)
			limit.TotalVolume -= order.size
			break
		}
	}
	order.limit = nil
	// resort the orders by timestamp
	sort.Slice(limit.orders, func(i, j int) bool {
		return limit.orders[i].timestamp < limit.orders[j].timestamp
	})
}

func (limit *Limit) matchOrder(order *Order) []Match {
	matches := []Match{}

	// iterate through each order in the limit
	for _, limitOrder := range limit.orders {
		// fill the order
		match := limit.fillOrders(limitOrder, order)
		// add the match to the list of matches
		matches = append(matches, match)

		limit.TotalVolume -= match.sizeFilled

		// remove the limit order if it is filled
		if limitOrder.size == 0 {
			limit.removeOrder(limitOrder)
		}

		// if the order is filled, break
		if order.size == 0 {
			break
		}
	}

	return matches
}

func (limit *Limit) fillOrders(limitOrder, order *Order) Match {
	var (
		ask        *Order
		bid        *Order
		sizeFilled uint64
	)

	if order.side == Bid {
		bid = order
		ask = limitOrder
	} else {
		bid = limitOrder
		ask = order
	}

	// if the order is smaller than the limit order, fill the order
	if limitOrder.size >= order.size {
		limitOrder.size -= order.size
		sizeFilled = order.size
		order.size = 0
	} else {
		// if the order is larger than the limit order, fill the limit order
		order.size -= limitOrder.size
		sizeFilled = limitOrder.size
		limitOrder.size = 0
	}

	return Match{
		ask:        ask,
		bid:        bid,
		sizeFilled: sizeFilled,
		price:      limit.Price,
	}
}
