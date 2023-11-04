package main

type Match struct {
	ask        *Order
	bid        *Order
	sizeFilled uint64
	price      uint64
}

type Limit struct {
	price       uint64
	orders      []*Order
	totalVolume uint64
}

func newLimit(price uint64) *Limit {
	return &Limit{
		price:  price,
		orders: []*Order{},
	}
}

func (limit *Limit) addOrder(order *Order) {
	order.limit = limit
	limit.orders = append(limit.orders, order)
	limit.totalVolume += order.size
}

func (limit *Limit) matchOrder(order *Order) []Match {
	matches := []Match{}

	// iterate through each order in the limit
	for _, limitOrder := range limit.orders {
		// fill the order
		match := limit.fillOrders(limitOrder, order)
		// add the match to the list of matches
		matches = append(matches, match)

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
		price:      limit.price,
	}
}
