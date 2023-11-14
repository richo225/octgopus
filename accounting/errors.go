package accounting

import "net/http"

type AccountNotFoundError struct {
	signer string
}

func (e *AccountNotFoundError) Error() string {
	return "AccountNotFound: " + e.signer
}

func (e *AccountNotFoundError) HTTPCode() int {
	return http.StatusNotFound
}

type AccountAlreadyExistsError struct {
	signer string
}

func (e *AccountAlreadyExistsError) Error() string {
	return "AccountAlreadyExists: " + e.signer
}

func (e *AccountAlreadyExistsError) HTTPCode() int {
	return http.StatusBadRequest
}

type AccountUnderFundedError struct {
	signer string
}

func (e *AccountUnderFundedError) Error() string {
	return "AccountUnderFunded: " + e.signer
}

func (e *AccountUnderFundedError) HTTPCode() int {
	return http.StatusBadRequest
}
