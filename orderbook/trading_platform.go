package orderbook

import (
	"encoding/json"
	"errors"

	"github.com/richo225/octgopus/accounting"
)

type OrderbookNotFoundError struct {
	pair TradingPair
}

func (e *OrderbookNotFoundError) Error() string {
	return "MarketNotfound : " + e.pair.ToString()
}

type OrderType string

const (
	LimitOrder  OrderType = "limit"
	MarketOrder OrderType = "market"
)

func (ot *OrderType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case string(LimitOrder), string(MarketOrder):
		*ot = OrderType(s)
	default:
		return errors.New("invalid order type")
	}

	return nil
}

type TradingPair struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

func NewTradingPair(base string, quote string) TradingPair {
	return TradingPair{
		Base:  base,
		Quote: quote,
	}
}

func (pair *TradingPair) ToString() string {
	return pair.Base + "/" + pair.Quote
}

type TradingPlatform struct {
	Accounts   *accounting.Accounts
	Orderbooks map[TradingPair]*Orderbook `json:"orderbooks"`
}

func NewTradingPlatform() *TradingPlatform {
	return &TradingPlatform{
		Accounts:   accounting.NewAccounts(),
		Orderbooks: make(map[TradingPair]*Orderbook),
	}
}

func (platform *TradingPlatform) AddNewMarket(pair TradingPair) *Orderbook {
	ob := newOrderBook()
	ob.Market = &pair
	platform.Orderbooks[pair] = ob

	return ob
}

func (platform *TradingPlatform) PlaceMarketOrder(pair TradingPair, order *Order) ([]Match, error) {
	orderbook, err := platform.GetOrderBook(pair)

	if err != nil {
		return nil, err
	}

	matches, err := orderbook.placeMarketOrder(order)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (platform *TradingPlatform) PlaceLimitOrder(pair TradingPair, price float64, order *Order) error {
	orderbook, err := platform.GetOrderBook(pair)

	if err != nil {
		return err
	}

	orderbook.placeLimitOrder(price, order)
	return nil
}

func (platform *TradingPlatform) GetOrderBook(pair TradingPair) (*Orderbook, error) {
	orderbook, ok := platform.Orderbooks[pair]
	if !ok {
		return nil, &OrderbookNotFoundError{pair}
	}

	return orderbook, nil
}
