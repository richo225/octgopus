package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
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
		return err
	}

	orderbook.GetAsks()
	orderbook.GetBids()

	return c.JSON(http.StatusOK, &orderbook)
}

// Orders
func (c *CustomContext) handleCreateOrder() error {
	params := PlaceOrderRequestParams{}
	c.Bind(&params)

	pair := orderbook.NewTradingPair(params.Base, params.Quote)
	order := orderbook.NewOrder(params.Side, params.Size)

	if params.Type == orderbook.MarketOrder {
		matches, err := c.platform.PlaceMarketOrder(pair, order)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &matches)
	}

	err := c.platform.PlaceLimitOrder(pair, params.Price, order)
	if err != nil {
		return err
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
		return err
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
		return err
	}

	return c.JSON(http.StatusOK, &tx)
}

func (c *CustomContext) handleAccountSend() error {
	params := AccountSendParams{}
	c.Bind(&params)

	tx, err := c.platform.Accounts.Send(params.Signer, params.Recipient, params.Amount)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &tx)
}
