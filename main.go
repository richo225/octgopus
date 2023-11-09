package main

import (
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	p := newTradingPlatform()

	pretty.Log("Seeding data...")
	p.seedData()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/", sayHello)

	e.GET("/orderbooks", p.handleGetOrderbook)
	e.POST("/orderbooks", p.handleCreateOrderbook)
	e.POST("/orders", p.handleCreateOrder)

	e.GET("/accounts", p.handleGetAccounts)
	e.GET("/accounts/:signer", p.handleGetAccountBalance)
	e.POST("/accounts/:signer", p.handleCreateAccount)
	e.POST("/accounts/:signer/deposit", p.handleAccountDeposit)
	e.POST("/accounts/:signer/withdraw", p.handleAccountWithdraw)
	e.POST("/accounts/:signer/send", p.handleAccountSend)

	pretty.Log("Starting server...")
	e.Logger.Fatal(e.Start(":8080"))
}
