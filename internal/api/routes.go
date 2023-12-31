package api

import (
	"github.com/labstack/echo/v4"
	"github.com/richo225/octgopus/internal/orderbook"
)

type CustomContext struct {
	echo.Context
	platform *orderbook.TradingPlatform
}

func registerHandlers(e *echo.Echo, p *orderbook.TradingPlatform) {
	withPlatform := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, p}
			return next(cc)
		}
	}

	e.GET("/", CheckHealth)

	orderbooks := e.Group("/orderbooks", withPlatform)
	orderbooks.GET("", withCustomContext((*CustomContext).handleGetOrderbook))
	orderbooks.GET("/reset", withCustomContext((*CustomContext).handleResetOrderbooks))
	orderbooks.POST("", withCustomContext((*CustomContext).handleCreateOrderbook))

	orders := e.Group("/orders", withPlatform)
	orders.POST("", withCustomContext((*CustomContext).handleCreateOrder))

	accounts := e.Group("/accounts", withPlatform)
	accounts.GET("", withCustomContext((*CustomContext).handleGetAccounts))
	accounts.GET("/:signer", withCustomContext((*CustomContext).handleGetAccountBalance))
	accounts.POST("/:signer", withCustomContext((*CustomContext).handleCreateAccount))
	accounts.POST("/:signer/deposit", withCustomContext((*CustomContext).handleAccountDeposit))
	accounts.POST("/:signer/withdraw", withCustomContext((*CustomContext).handleAccountWithdraw))
	accounts.POST("/:signer/send", withCustomContext((*CustomContext).handleAccountSend))
}

func withCustomContext(handler func(c *CustomContext) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c.(*CustomContext))
	}
}
