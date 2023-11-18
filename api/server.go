package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/richo225/octgopus/internal/orderbook"
)

func Start(p *orderbook.TradingPlatform) {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.Binder = &Binder{}
	e.Validator = NewValidator()
	e.HTTPErrorHandler = HTTPErrorHandler

	registerHandlers(e, p)

	pretty.Log("Starting server...")

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
