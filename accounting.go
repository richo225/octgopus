package main

type AccountNotFoundError struct {
	message string
}

func (e *AccountNotFoundError) Error() string {
	return "AccountNotFound : " + e.message
}

type AccountUnderFundedError struct {
	message string
}

func (e *AccountUnderFundedError) Error() string {
	return "AccountUnderFunded : " + e.message
}

type Accounts struct {
	// Stores the total balance of each account.
	accounts map[string]uint64
}

func newAccounts() *Accounts {
	return &Accounts{
		accounts: make(map[string]uint64),
	}
}

func (a *Accounts) balanceOf(signer string) (uint64, error) {
	balance, ok := a.accounts[signer]
	if !ok {
		return 0, &AccountNotFoundError{signer}
	}

	return balance, nil
}

func (a *Accounts) deposit(signer string, amount uint64) error {
	balance, err := a.balanceOf(signer)
	if err != nil {
		return err
	}

	newBalance := balance + amount
	a.accounts[signer] = newBalance

	return nil
}

func (a *Accounts) withdraw(signer string, amount uint64) error {
	balance, err := a.balanceOf(signer)
	if err != nil {
		return &AccountNotFoundError{"account not found"}
	} else {
		if balance < amount {
			return &AccountUnderFundedError{signer}
		}
		newBalance := balance - amount
		a.accounts[signer] = newBalance
	}

	return nil
}