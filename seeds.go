package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func (platform *TradingPlatform) seedData() {
	ethusd := newTradingPair("ETH", "USD")
	ethgbp := newTradingPair("ETH", "GBP")
	btcusd := newTradingPair("BTC", "USD")
	btcgbp := newTradingPair("BTC", "GBP")

	platform.writeFromFile("data/eth_usd_order_book.csv", ethusd)
	platform.writeFromFile("data/eth_gbp_order_book.csv", ethgbp)
	platform.writeFromFile("data/btc_usd_order_book.csv", btcusd)
	platform.writeFromFile("data/btc_gbp_order_book.csv", btcgbp)
}

func (platform *TradingPlatform) writeFromFile(filepath string, pair TradingPair) {
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

	platform.addNewMarket(pair)

	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		askPrice, _ := strconv.ParseFloat(record[0], 64)
		askAmount, _ := strconv.ParseFloat(record[1], 64)
		bidPrice, _ := strconv.ParseFloat(record[2], 64)
		bidAmount, _ := strconv.ParseFloat(record[3], 64)

		askOrder := newOrder(Ask, askAmount)
		platform.placeLimitOrder(pair, askPrice, askOrder)

		bidOrder := newOrder(Bid, bidAmount)
		platform.placeLimitOrder(pair, bidPrice, bidOrder)
	}
}
