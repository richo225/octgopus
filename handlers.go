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
	return c.String(http.StatusOK, "Order")
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
