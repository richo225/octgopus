package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Binder struct{}

func (cb *Binder) Bind(i interface{}, c echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			he.Code = http.StatusBadRequest
		}
	}

	return
}
