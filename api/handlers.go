package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/richo225/octgopus/accounting"
	"github.com/richo225/octgopus/orderbook"
)

type CustomContext struct {
	echo.Context
	platform *orderbook.TradingPlatform
}

func CheckHealth(c echo.Context) error {
	return c.String(http.StatusOK, "Up")
}

// Orderbooks
func (c *CustomContext) handleCreateOrderbook() error {
	params := MarketParams{}
	c.Bind(&params)
	pair := orderbook.NewTradingPair(params.Base, params.Quote)

	orderbook := c.platform.AddNewMarket(pair)

	return c.JSON(http.StatusOK, &orderbook)
}

func (c *CustomContext) handleGetOrderbook() error {
	params := MarketParams{}
	c.Bind(&params)
	pair := orderbook.NewTradingPair(params.Base, params.Quote)

	orderbook, err := c.platform.GetOrderBook(pair)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	orderbook.GetAsks()
	orderbook.GetBids()

	return c.JSON(http.StatusOK, &orderbook)
}

// Orders
func (c *CustomContext) handleCreateOrder() error {
	params := PlaceOrderRequestParams{}
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pair := orderbook.NewTradingPair(params.Base, params.Quote)
	order := orderbook.NewOrder(params.Side, params.Size)

	if params.Type == orderbook.MarketOrder {
		matches, err := c.platform.PlaceMarketOrder(pair, order)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, &matches)
	}

	err := c.platform.PlaceLimitOrder(pair, params.Price, order)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &order)
}

// Accounting
func (c *CustomContext) handleCreateAccount() error {
	signer := c.Param("signer")

	err := c.platform.Accounts.CreateAccount(signer)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Account created")
}

func (c *CustomContext) handleGetAccounts() error {
	accounts := c.platform.Accounts

	return c.JSON(http.StatusOK, &accounts)
}

func (c *CustomContext) handleGetAccountBalance() error {
	params := AccountBalanceParams{}
	c.Bind(&params)

	balance, err := c.platform.Accounts.BalanceOf(params.Signer)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("%f", balance))
}

func (c *CustomContext) handleAccountDeposit() error {
	params := AccountActionParams{}
	c.Bind(&params)

	tx := c.platform.Accounts.Deposit(params.Signer, params.Amount)
	return c.JSON(http.StatusOK, &tx)
}

func (c *CustomContext) handleAccountWithdraw() error {
	params := AccountActionParams{}
	c.Bind(&params)

	tx, err := c.platform.Accounts.Withdraw(params.Signer, params.Amount)
	if err != nil {
		if _, ok := err.(*accounting.AccountNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if _, ok := err.(*accounting.AccountUnderFundedError); ok {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	return c.JSON(http.StatusOK, &tx)
}

func (c *CustomContext) handleAccountSend() error {
	params := AccountSendParams{}
	c.Bind(&params)

	tx, err := c.platform.Accounts.Send(params.Signer, params.Recipient, params.Amount)
	if err != nil {
		if _, ok := err.(*accounting.AccountNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if _, ok := err.(*accounting.AccountUnderFundedError); ok {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	return c.JSON(http.StatusOK, &tx)
}

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
