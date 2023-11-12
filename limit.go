package main

import "sort"

type Match struct {
	Ask        *Order  `json:"ask"`
	Bid        *Order  `json:"bid"`
	SizeFilled float64 `json:"size_filled"`
	Price      float64 `json:"price"`
}

type Limit struct {
	Price       float64  `json:"price"`
	TotalVolume float64  `json:"total_volume"`
	Orders      []*Order `json:"orders"`
}

func newLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func (limit *Limit) addOrder(order *Order) {
	order.Price = limit.Price
	limit.Orders = append(limit.Orders, order)
	limit.TotalVolume += order.Size
}

func (limit *Limit) removeOrder(order *Order) {
	for i, o := range limit.Orders {
		if o == order {
			limit.Orders = append(limit.Orders[:i], limit.Orders[i+1:]...)
			limit.TotalVolume -= order.Size
			break
		}
	}
	order.Price = 0
	// resort the orders by timestamp
	sort.Slice(limit.Orders, func(i, j int) bool {
		return limit.Orders[i].Timestamp < limit.Orders[j].Timestamp
	})
}

func (limit *Limit) matchOrder(order *Order) []Match {
	matches := []Match{}

	// iterate through each order in the limit
	for _, limitOrder := range limit.Orders {
		// fill the order
		match := limit.fillOrders(limitOrder, order)
		// add the match to the list of matches
		matches = append(matches, match)

		limit.TotalVolume -= match.SizeFilled

		// remove the limit order if it is filled
		if limitOrder.Size == 0 {
			limit.removeOrder(limitOrder)
		}

		// if the order is filled, break
		if order.Size == 0 {
			break
		}
	}

	return matches
}

func (limit *Limit) fillOrders(limitOrder, order *Order) Match {
	var (
		ask        *Order
		bid        *Order
		sizeFilled float64
	)

	if order.Side == Bid {
		bid = order
		ask = limitOrder
	} else {
		bid = limitOrder
		ask = order
	}

	// if the order is smaller than the limit order, fill the order
	if limitOrder.Size >= order.Size {
		limitOrder.Size -= order.Size
		sizeFilled = order.Size
		order.Size = 0
	} else {
		// if the order is larger than the limit order, fill the limit order
		order.Size -= limitOrder.Size
		sizeFilled = limitOrder.Size
		limitOrder.Size = 0
	}

	return Match{
		Ask:        ask,
		Bid:        bid,
		SizeFilled: sizeFilled,
		Price:      limit.Price,
	}
}
