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
	// Open the CSV file
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a new reader
	r := csv.NewReader(f)
	// Read the header row
	_, err = r.Read()
	if err != nil {
		panic(err)
	}

	// Create market for pair
	platform.addNewMarket(pair)

	// Iterate over the remaining rows
	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		// Parse price and amount
		askPrice, _ := strconv.ParseUint(record[0], 10, 64)
		askAmount, _ := strconv.ParseUint(record[1], 10, 64)
		bidPrice, _ := strconv.ParseUint(record[2], 10, 64)
		bidAmount, _ := strconv.ParseUint(record[3], 10, 64)

		// Create and place ask order
		askOrder := newOrder(Ask, askAmount)
		platform.placeLimitOrder(pair, askPrice, askOrder)

		// Create and place bid order
		bidOrder := newOrder(Bid, bidAmount)
		platform.placeLimitOrder(pair, bidPrice, bidOrder)
	}
}
