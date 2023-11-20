package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return nil
}
