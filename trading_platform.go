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
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

func newTradingPair(base string, quote string) TradingPair {
	return TradingPair{
		Base:  base,
		Quote: quote,
	}
}

func (pair *TradingPair) toString() string {
	return pair.Base + "/" + pair.Quote
}

type TradingPlatform struct {
	accounts   *Accounts
	Orderbooks map[TradingPair]*Orderbook `json:"orderbooks"`
}

func newTradingPlatform() *TradingPlatform {
	return &TradingPlatform{
		accounts:   newAccounts(),
		Orderbooks: make(map[TradingPair]*Orderbook),
	}
}

func (platform *TradingPlatform) addNewMarket(pair TradingPair) *Orderbook {
	ob := newOrderBook()
	ob.Market = &pair
	platform.Orderbooks[pair] = ob

	return ob
}

func (platform *TradingPlatform) placeMarketOrder(pair TradingPair, order *Order) ([]Match, error) {
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
	orderbook, ok := platform.Orderbooks[pair]
	if !ok {
		return nil, &OrderbookNotFoundError{pair}
	}

	return orderbook, nil
}
