package accounting

type AccountNotFoundError struct {
	signer string
}

func (e *AccountNotFoundError) Error() string {
	return "AccountNotFound: " + e.signer
}

type AccountAlreadyExistsError struct {
	signer string
}

func (e *AccountAlreadyExistsError) Error() string {
	return "AccountAlreadyExists: " + e.signer
}

type AccountUnderFundedError struct {
	signer string
}

func (e *AccountUnderFundedError) Error() string {
	return "AccountUnderFunded: " + e.signer
}
