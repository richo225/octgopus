package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func checkHealth(c echo.Context) error {
	return c.String(http.StatusOK, "Up")
}

// Orderbooks
func (platform *TradingPlatform) handleCreateOrderbook(c echo.Context) error {
	params := MarketParams{}
	c.Bind(&params)
	pair := newTradingPair(params.Base, params.Quote)

	orderbook := platform.addNewMarket(pair)

	return c.JSON(http.StatusOK, &orderbook)
}

func (platform *TradingPlatform) handleGetOrderbook(c echo.Context) error {
	params := MarketParams{}
	c.Bind(&params)
	pair := newTradingPair(params.Base, params.Quote)

	orderbook, err := platform.getOrderBook(pair)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	orderbook.getAsks()
	orderbook.getBids()

	return c.JSON(http.StatusOK, &orderbook)
}

// Orders
func (platform *TradingPlatform) handleCreateOrder(c echo.Context) error {
	base := c.QueryParam("base")
	quote := c.QueryParam("quote")
	pair := newTradingPair(base, quote)
	side := c.QueryParam("side")
	orderType := c.QueryParam("type")

	// TODO: replace above with binding to a request struct

	priceStr := c.QueryParam("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sizeStr := c.QueryParam("size")
	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order := newOrder(Side(side), size)
	if OrderType(orderType) == MarketOrder {
		matches, err := platform.placeMarketOrder(pair, order)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, &matches)
	}

	err = platform.placeLimitOrder(pair, price, order)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &order)
}

// Accounting
func (platform *TradingPlatform) handleCreateAccount(c echo.Context) error {
	signer := c.Param("signer")

	err := platform.accounts.createAccount(signer)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Account created")
}

func (platform *TradingPlatform) handleGetAccounts(c echo.Context) error {
	accounts := platform.accounts

	return c.JSON(http.StatusOK, &accounts)
}

func (platform *TradingPlatform) handleGetAccountBalance(c echo.Context) error {
	params := AccountBalanceParams{}
	c.Bind(&params)

	balance, err := platform.accounts.balanceOf(params.Signer)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("%f", balance))
}

func (platform *TradingPlatform) handleAccountDeposit(c echo.Context) error {
	params := AccountActionParams{}
	c.Bind(&params)

	tx := platform.accounts.deposit(params.Signer, params.Amount)
	return c.JSON(http.StatusOK, &tx)
}

func (platform *TradingPlatform) handleAccountWithdraw(c echo.Context) error {
	params := AccountActionParams{}
	c.Bind(&params)

	tx, err := platform.accounts.withdraw(params.Signer, params.Amount)
	if err != nil {
		if _, ok := err.(*AccountNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if _, ok := err.(*AccountUnderFundedError); ok {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	return c.JSON(http.StatusOK, &tx)
}

func (platform *TradingPlatform) handleAccountSend(c echo.Context) error {
	params := AccountSendParams{}
	c.Bind(&params)

	tx, err := platform.accounts.send(params.Signer, params.Recipient, params.Amount)
	if err != nil {
		if _, ok := err.(*AccountNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if _, ok := err.(*AccountUnderFundedError); ok {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	return c.JSON(http.StatusOK, &tx)
}

type PlaceOrderRequestParams struct {
	*MarketParams
	Side      Side      `json:"side" form:"side" query:"side"`
	OrderType OrderType `json:"order_type" form:"order_type" query:"order_type"`
	Price     float64   `json:"price" form:"price" query:"price"`
	Size      float64   `json:"size" form:"size" query:"size"`
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
	*AccountActionParams
	Recipient string `json:"recipient" form:"recipient" query:"recipient"`
}
