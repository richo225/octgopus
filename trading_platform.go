package main

type OrderbookNotFoundError struct {
	pair TradingPair
}

func (e *OrderbookNotFoundError) Error() string {
	return "MarketNotfound : " + e.pair.toString()
}

type OrderType string

const (
	LimitOrder  OrderType = "limit"
	MarketOrder OrderType = "market"
)

type TradingPair struct {
	base  string // BTC
	quote string // USD
}

func newTradingPair(base string, quote string) TradingPair {
	return TradingPair{
		base:  base,
		quote: quote,
	}
}

func (pair *TradingPair) toString() string {
	return pair.base + "/" + pair.quote
}

type TradingPlatform struct {
	orderbooks map[TradingPair]Orderbook
}

func newTradingPlatform() *TradingPlatform {
	return &TradingPlatform{
		orderbooks: make(map[TradingPair]Orderbook),
	}
}

func (platform *TradingPlatform) addNewMarket(pair TradingPair) {
	platform.orderbooks[pair] = *newOrderBook()
}

func (platform *TradingPlatform) placeOrder(pair TradingPair, price uint64, side Side, size uint64, orderType OrderType) ([]Match, error) {

	order := newOrder(side, size)

	if orderType == MarketOrder {
		return platform.placeMarketOrder(pair, price, order)
	} else {
		return nil, platform.placeLimitOrder(pair, price, order)
	}
}

func (platform *TradingPlatform) placeMarketOrder(pair TradingPair, price uint64, order *Order) ([]Match, error) {
	orderbook, err := platform.getOrderBook(pair)

	if err != nil {
		return nil, err
	}

	matches, err := orderbook.placeMarketOrder(order)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (platform *TradingPlatform) placeLimitOrder(pair TradingPair, price uint64, order *Order) error {
	orderbook, err := platform.getOrderBook(pair)

	if err != nil {
		return err
	}

	orderbook.placeLimitOrder(price, order)
	return nil
}

func (platform *TradingPlatform) getOrderBook(pair TradingPair) (*Orderbook, error) {
	// check in orderbooks for the trading pair
	orderbook, ok := platform.orderbooks[pair]
	if !ok {
		// if market does not exist, return error
		return nil, &OrderbookNotFoundError{pair}
	}

	return &orderbook, nil
}
