package api

import "github.com/richo225/octgopus/orderbook"

type PlaceOrderRequestParams struct {
	MarketParams
	Side  orderbook.Side      `json:"side" form:"side" query:"side"`
	Type  orderbook.OrderType `json:"type" form:"type" query:"type"`
	Price float64             `json:"price" form:"price" query:"price"`
	Size  float64             `json:"size" form:"size" query:"size"`
}

type MarketParams struct {
	Quote string `json:"quote" form:"quote" query:"quote"`
	Base  string `json:"base" form:"base" query:"base"`
}

type AccountBalanceParams struct {
	Signer string `json:"signer" form:"signer" query:"signer"`
}

type AccountActionParams struct {
	Signer string  `json:"signer" form:"signer" query:"signer"`
	Amount float64 `json:"amount" form:"amount" query:"amount"`
}
type AccountSendParams struct {
	AccountActionParams
	Recipient string `json:"recipient" form:"recipient" query:"recipient"`
}
