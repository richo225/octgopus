package main

type Side string

const (
	Bid Side = "bid"
	Ask Side = "ask"
)

type Order struct {
	side Side
	size uint64
}

func newOrder(side Side, size uint64) Order {
	return Order{side, size}
}

type Limit struct {
	price  uint64
	orders []Order
}

func newLimit(price uint64) *Limit {
	return &Limit{
		price:  price,
		orders: make([]Order, 0),
	}
}

func (limit *Limit) addOrder(order Order) {
	limit.orders = append(limit.orders, order)
}

type Orderbook struct {
	asks map[uint64]Limit
	bids map[uint64]Limit
}

func newOrderBook() *Orderbook {
	return &Orderbook{
		asks: make(map[uint64]Limit),
		bids: make(map[uint64]Limit),
	}
}

func (book *Orderbook) addOrder(price uint64, order Order) {
	if order.side == Bid {
		// check if a limit already exists at the order price
		limit, ok := book.bids[price]
		if ok {
			// if one exists, add the order to that limit
			limit.addOrder(order)
		} else {
			// if doesn't exist, create a new limit with the order
			newLimit := newLimit(price)
			newLimit.addOrder(order)
			book.bids[price] = *newLimit
		}
	} else {
		limit, ok := book.asks[price]
		if ok {
			limit.addOrder(order)
		} else {
			newLimit := newLimit(price)
			newLimit.addOrder(order)
			book.asks[price] = *newLimit
		}
	}
}
