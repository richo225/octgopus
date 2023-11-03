package main

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

func (limit *Limit) matchOrder(order *Order) {
	// iterate through each order in the limit
	for _, limitOrder := range limit.orders {
		// fill the order
		limit.fillOrders(limitOrder, order)

		if order.size == 0 {
			break
		}
	}
}

func (limit *Limit) fillOrders(limitOrder, order *Order) {
	// if the order is smaller than the limit order, fill the order
	if limitOrder.size >= order.size {
		limitOrder.size -= order.size
		order.size = 0
	} else {
		// if the order is larger than the limit order, fill the limit order
		order.size -= limitOrder.size
		limitOrder.size = 0
	}
}
