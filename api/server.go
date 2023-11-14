package api

import (
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/richo225/octgopus/orderbook"
)

func Start(p *orderbook.TradingPlatform) {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/", CheckHealth)

	withPlatform := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, p}
			return next(cc)
		}
	}

	orderbooks := e.Group("/orderbooks", withPlatform)
	orderbooks.GET("", withCustomContext((*CustomContext).handleGetOrderbook))
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

	pretty.Log("Starting server...")
	e.Logger.Fatal(e.Start("localhost:8080"))
}

func withCustomContext(handler func(c *CustomContext) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c.(*CustomContext))
	}
}
