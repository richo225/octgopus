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
