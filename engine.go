package main

type OrderbookNotFoundError struct {
	pair TradingPair
}

func (e *OrderbookNotFoundError) Error() string {
	return "MarketNotfound : " + e.pair.toString()
}

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

type MatchingEngine struct {
	orderbooks map[TradingPair]Orderbook
}

func newMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		orderbooks: make(map[TradingPair]Orderbook),
	}
}

func (engine *MatchingEngine) addNewMarket(pair TradingPair) {
	engine.orderbooks[pair] = *newOrderBook()
}

func (engine *MatchingEngine) placeMarketOrder(pair TradingPair, price uint64, order *Order) ([]Match, error) {
	orderbook, err := engine.getOrderBook(pair)

	if err != nil {
		return nil, err
	}

	matches, err := orderbook.placeMarketOrder(order)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (engine *MatchingEngine) placeLimitOrder(pair TradingPair, price uint64, order *Order) error {
	orderbook, err := engine.getOrderBook(pair)

	if err != nil {
		return err
	}

	orderbook.placeLimitOrder(price, order)
	return nil
}

func (engine *MatchingEngine) getOrderBook(pair TradingPair) (*Orderbook, error) {
	// check in orderbooks for the trading pair
	orderbook, ok := engine.orderbooks[pair]
	if !ok {
		// if market does not exist, return error
		return nil, &OrderbookNotFoundError{pair}
	}

	return &orderbook, nil
}
