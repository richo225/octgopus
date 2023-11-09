package main

func (platform *TradingPlatform) seedData() {
	platform.accounts.deposit("Alice", 25000)
	platform.accounts.deposit("Bob", 460)
	platform.accounts.deposit("Charlie", 2690)

	ethusd := newTradingPair("ETH", "USD")
	ethgp := newTradingPair("ETH", "GBP")
	btcusd := newTradingPair("BTC", "USD")
	btcgbp := newTradingPair("BTC", "GBP")

	platform.addNewMarket(ethusd)
	platform.addNewMarket(ethgp)
	platform.addNewMarket(btcusd)
	platform.addNewMarket(btcgbp)

	platform.placeLimitOrder(ethusd, 1680, newOrder(Bid, 2))
	platform.placeLimitOrder(ethusd, 1640, newOrder(Bid, 5))
	platform.placeLimitOrder(ethusd, 1620, newOrder(Bid, 5))
	platform.placeLimitOrder(ethusd, 1700, newOrder(Ask, 1))
	platform.placeLimitOrder(ethusd, 1720, newOrder(Ask, 4))
	platform.placeLimitOrder(ethusd, 1740, newOrder(Ask, 12))

	platform.placeLimitOrder(btcusd, 1680, newOrder(Bid, 2))
	platform.placeLimitOrder(btcusd, 1640, newOrder(Bid, 5))
	platform.placeLimitOrder(btcusd, 1620, newOrder(Bid, 5))
	platform.placeLimitOrder(btcusd, 1700, newOrder(Ask, 1))
	platform.placeLimitOrder(btcusd, 1720, newOrder(Ask, 4))
	platform.placeLimitOrder(btcusd, 1740, newOrder(Ask, 12))
}
