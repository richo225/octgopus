package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPError interface {
	HTTPCode() int
}

func HTTPErrorHandler(err error, c echo.Context) {
	type Response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	code := http.StatusInternalServerError
	message := err.Error()

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	} else if he, ok := err.(HTTPError); ok {
		code = he.HTTPCode()
	}

	c.JSON(code, &Response{
		Code:    code,
		Message: message,
	})
}
