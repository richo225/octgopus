package orderbook

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"

	"github.com/kr/pretty"
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

func (platform *TradingPlatform) Reset() {
	platform.Orderbooks = make(map[TradingPair]*Orderbook)
	platform.Accounts = accounting.NewAccounts()

	platform.SeedData()
}

func (platform *TradingPlatform) SeedData() {
	pretty.Log("Seeding data...")

	ethusd := NewTradingPair("ETH", "USD")
	ethgbp := NewTradingPair("ETH", "GBP")
	btcusd := NewTradingPair("BTC", "USD")
	btcgbp := NewTradingPair("BTC", "GBP")

	writeFromFile(platform, "data/eth_usd_order_book.csv", ethusd)
	writeFromFile(platform, "data/eth_gbp_order_book.csv", ethgbp)
	writeFromFile(platform, "data/btc_usd_order_book.csv", btcusd)
	writeFromFile(platform, "data/btc_gbp_order_book.csv", btcgbp)

	pretty.Log("Seeding data complete!")
}

func writeFromFile(platform *TradingPlatform, filepath string, pair TradingPair) {
	pretty.Log("Seeding data for " + pair.ToString() + "...")

	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil {
		panic(err)
	}

	platform.AddNewMarket(pair)

	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		askPrice, _ := strconv.ParseFloat(record[0], 64)
		askAmount, _ := strconv.ParseFloat(record[1], 64)
		bidPrice, _ := strconv.ParseFloat(record[2], 64)
		bidAmount, _ := strconv.ParseFloat(record[3], 64)

		askOrder := NewOrder(Ask, askAmount)
		platform.PlaceLimitOrder(pair, askPrice, askOrder)

		bidOrder := NewOrder(Bid, bidAmount)
		platform.PlaceLimitOrder(pair, bidPrice, bidOrder)
	}

	pretty.Log("Seeding data for " + pair.ToString() + " complete!")
}
