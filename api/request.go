package api

import "github.com/richo225/octgopus/orderbook"

type PlaceOrderRequestParams struct {
	MarketParams
	Side  orderbook.Side      `json:"side" form:"side" query:"side" validate:"required"`
	Type  orderbook.OrderType `json:"type" form:"type" query:"type" validate:"required"`
	Price float64             `json:"price" form:"price" query:"price" validate:"required"`
	Size  float64             `json:"size" form:"size" query:"size" validate:"required"`
}

type MarketParams struct {
	Quote string `json:"quote" form:"quote" query:"quote" validate:"required"`
	Base  string `json:"base" form:"base" query:"base" validate:"required"`
}

type AccountBalanceParams struct {
	Signer string `json:"signer" form:"signer" query:"signer" validate:"required"`
}

type AccountActionParams struct {
	Signer string  `json:"signer" form:"signer" query:"signer" validate:"required"`
	Amount float64 `json:"amount" form:"amount" query:"amount" validate:"required"`
}
type AccountSendParams struct {
	AccountActionParams
	Recipient string `json:"recipient" form:"recipient" query:"recipient" validate:"required"`
}
