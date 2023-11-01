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
	if ok {
		return balance, nil
	} else {
		return 0, &AccountNotFoundError{signer}
	}
}

// Either deposits the `amount` provided into the `signer` account or adds the amount to the existing account
func (a *Accounts) deposit(signer string, amount uint64) error {
	balance, err := a.balanceOf(signer)
	if err != nil {
		a.accounts[signer] = amount
	} else {
		newBalance := balance + amount
		a.accounts[signer] = newBalance
	}

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

func (a *Accounts) send(sender string, recipient string, amount uint64) error {
	err := a.withdraw(sender, amount)
	if err != nil {
		return err
	}

	err = a.deposit(recipient, amount)
	if err != nil {
		return err
	}

	return nil
}
