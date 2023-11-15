package orderbook

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/richo225/octgopus/internal/accounting"
)

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

	mu sync.RWMutex
}

func NewTradingPlatform() *TradingPlatform {
	return &TradingPlatform{
		Accounts:   accounting.NewAccounts(),
		Orderbooks: make(map[TradingPair]*Orderbook),
	}
}

func (platform *TradingPlatform) AddNewMarket(pair TradingPair) *Orderbook {
	platform.mu.Lock()
	defer platform.mu.Unlock()

	ob := newOrderBook()
	ob.Market = &pair
	platform.Orderbooks[pair] = ob

	return ob
}

func (platform *TradingPlatform) PlaceMarketOrder(pair TradingPair, order *Order) ([]Match, error) {
	platform.mu.RLock()
	orderbook, err := platform.GetOrderBook(pair)
	platform.mu.RUnlock()

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
	platform.mu.RLock()
	orderbook, err := platform.GetOrderBook(pair)
	platform.mu.RUnlock()

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
