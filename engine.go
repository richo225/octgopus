package main

type OrderbookNotFoundError struct {
	message string
}

func (e *OrderbookNotFoundError) Error() string {
	return "MarketNotfound : " + e.message
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

func (engine *MatchingEngine) placeLimitOrder(pair TradingPair, price uint64, order *Order) error {
	// check in orderbooks for the trading pair
	orderbook, ok := engine.orderbooks[pair]
	if !ok {
		// if market does not exist, return error
		return &OrderbookNotFoundError{pair.toString()}
	}
	// if market already exists, add order to orderbook
	orderbook.addOrder(price, order)
	return nil
}
