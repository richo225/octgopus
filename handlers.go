package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func sayHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// Orderbooks
func (platform *TradingPlatform) handleCreateOrderbook(c echo.Context) error {
	base := c.QueryParam("base")
	quote := c.QueryParam("quote")
	pair := newTradingPair(base, quote)

	orderbook := platform.addNewMarket(pair)

	return c.JSON(http.StatusOK, &orderbook)
}

func (platform *TradingPlatform) handleGetOrderbook(c echo.Context) error {
	base := c.QueryParam("base")
	quote := c.QueryParam("quote")
	pair := newTradingPair(base, quote)

	orderbook, err := platform.getOrderBook(pair)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

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
	price, err := strconv.ParseUint(priceStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sizeStr := c.QueryParam("size")
	size, err := strconv.ParseUint(sizeStr, 10, 64)
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

func (platform *TradingPlatform) handleGetAccountBalance(c echo.Context) error {
	signer := c.Param("signer")
	balance, err := platform.accounts.balanceOf(signer)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("%d", balance))
}

func (platform *TradingPlatform) handleAccountDeposit(c echo.Context) error {
	signer := c.Param("signer")
	amountStr := c.QueryParam("amount")
	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx, err := platform.accounts.deposit(signer, amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, &tx)
}

func (platform *TradingPlatform) handleAccountWithdraw(c echo.Context) error {
	signer := c.Param("signer")
	amountStr := c.QueryParam("amount")
	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx, err := platform.accounts.withdraw(signer, amount)
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
	signer := c.Param("signer")
	recipient := c.QueryParam("recipient")
	amountStr := c.QueryParam("amount")
	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx, err := platform.accounts.send(signer, recipient, amount)
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
	Quote     string    `json:"quote" form:"quote" query:"quote"`
	Base      string    `json:"base" form:"base" query:"base"`
	Side      Side      `json:"side" form:"side" query:"side"`
	OrderType OrderType `json:"order_type" form:"order_type" query:"order_type"`
	Price     uint64    `json:"price" form:"price" query:"price"`
	Size      uint64    `json:"size" form:"size" query:"size"`
}
