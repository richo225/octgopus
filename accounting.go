package main

type TxAction string

const (
	Deposit  TxAction = "deposit"
	Withdraw TxAction = "withdraw"
)

type Tx struct {
	Action TxAction `json:"action"`
	Signer string   `json:"signer"`
	Amount uint64   `json:"amount"`
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

func (a *Accounts) createAccount(signer string) error {
	_, ok := a.accounts[signer]
	if ok {
		return &AccountAlreadyExistsError{signer}
	} else {
		a.accounts[signer] = 0
		return nil
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

func (a *Accounts) deposit(signer string, amount uint64) *Tx {
	var newBalance uint64

	balance, err := a.balanceOf(signer)
	if err != nil {
		newBalance = amount
	} else {
		newBalance = balance + amount
	}

	a.accounts[signer] = newBalance

	return &Tx{
		Action: Deposit,
		Signer: signer,
		Amount: amount,
	}
}

func (a *Accounts) withdraw(signer string, amount uint64) (*Tx, error) {
	balance, err := a.balanceOf(signer)
	if err != nil {
		return nil, &AccountNotFoundError{signer}
	} else {
		if balance < amount {
			return nil, &AccountUnderFundedError{signer}
		}
		newBalance := balance - amount
		a.accounts[signer] = newBalance

		return &Tx{
			Action: Withdraw,
			Signer: signer,
			Amount: amount,
		}, nil
	}
}

func (a *Accounts) send(sender string, recipient string, amount uint64) ([]*Tx, error) {
	wtx, err := a.withdraw(sender, amount)
	if err != nil {
		return nil, err
	}

	dtx := a.deposit(recipient, amount)

	return []*Tx{wtx, dtx}, nil
}
