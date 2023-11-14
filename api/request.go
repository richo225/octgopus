package api

import "github.com/richo225/octgopus/orderbook"

type PlaceOrderRequestParams struct {
	MarketParams
	Side  orderbook.Side      `json:"side" form:"side" query:"side" required:"true"`
	Type  orderbook.OrderType `json:"type" form:"type" query:"type" required:"true"`
	Price float64             `json:"price" form:"price" query:"price" required:"true"`
	Size  float64             `json:"size" form:"size" query:"size" required:"true"`
}

type MarketParams struct {
	Quote string `json:"quote" form:"quote" query:"quote" required:"true"`
	Base  string `json:"base" form:"base" query:"base" required:"true"`
}

type AccountBalanceParams struct {
	Signer string `json:"signer" form:"signer" query:"signer" required:"true"`
}

type AccountActionParams struct {
	Signer string  `json:"signer" form:"signer" query:"signer" required:"true"`
	Amount float64 `json:"amount" form:"amount" query:"amount" required:"true"`
}
type AccountSendParams struct {
	AccountActionParams
	Recipient string `json:"recipient" form:"recipient" query:"recipient" required:"true"`
}
