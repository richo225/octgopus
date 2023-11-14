package data

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/kr/pretty"
	"github.com/richo225/octgopus/orderbook"
)

func SeedOrderBook(platform *orderbook.TradingPlatform) {
	pretty.Log("Seeding data...")

	ethusd := orderbook.NewTradingPair("ETH", "USD")
	ethgbp := orderbook.NewTradingPair("ETH", "GBP")
	btcusd := orderbook.NewTradingPair("BTC", "USD")
	btcgbp := orderbook.NewTradingPair("BTC", "GBP")

	writeFromFile(platform, "data/eth_usd_order_book.csv", ethusd)
	writeFromFile(platform, "data/eth_gbp_order_book.csv", ethgbp)
	writeFromFile(platform, "data/btc_usd_order_book.csv", btcusd)
	writeFromFile(platform, "data/btc_gbp_order_book.csv", btcgbp)

	pretty.Log("Seeding data complete!")
}

func writeFromFile(platform *orderbook.TradingPlatform, filepath string, pair orderbook.TradingPair) {
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

		askOrder := orderbook.NewOrder(orderbook.Ask, askAmount)
		platform.PlaceLimitOrder(pair, askPrice, askOrder)

		bidOrder := orderbook.NewOrder(orderbook.Bid, bidAmount)
		platform.PlaceLimitOrder(pair, bidPrice, bidOrder)
	}

	pretty.Log("Seeding data for " + pair.ToString() + " complete!")
}
