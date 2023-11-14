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

	e.Binder = &Binder{}
	e.Validator = NewValidator()
	e.HTTPErrorHandler = HTTPErrorHandler

	registerHandlers(e, p)

	pretty.Log("Starting server...")
	e.Logger.Fatal(e.Start("localhost:8080"))
}
