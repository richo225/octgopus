package accounting

type TxAction string

const (
	Deposit  TxAction = "deposit"
	Withdraw TxAction = "withdraw"
)

type Tx struct {
	Action TxAction `json:"action"`
	Signer string   `json:"signer"`
	Amount float64  `json:"amount"`
}

type Accounts struct {
	// Stores the total balance of each account.
	Accounts map[string]float64 `json:"accounts"`
}

func NewAccounts() *Accounts {
	return &Accounts{
		Accounts: make(map[string]float64),
	}
}

func (a *Accounts) CreateAccount(signer string) error {
	_, ok := a.Accounts[signer]
	if ok {
		return &AccountAlreadyExistsError{signer}
	} else {
		a.Accounts[signer] = 0
		return nil
	}
}

func (a *Accounts) BalanceOf(signer string) (float64, error) {
	balance, ok := a.Accounts[signer]
	if ok {
		return balance, nil
	} else {
		return 0, &AccountNotFoundError{signer}
	}
}

func (a *Accounts) Deposit(signer string, amount float64) *Tx {
	var newBalance float64

	balance, err := a.BalanceOf(signer)
	if err != nil {
		newBalance = amount
	} else {
		newBalance = balance + amount
	}

	a.Accounts[signer] = newBalance

	return &Tx{
		Action: Deposit,
		Signer: signer,
		Amount: amount,
	}
}

func (a *Accounts) Withdraw(signer string, amount float64) (*Tx, error) {
	balance, err := a.BalanceOf(signer)
	if err != nil {
		return nil, &AccountNotFoundError{signer}
	} else {
		if balance < amount {
			return nil, &AccountUnderFundedError{signer}
		}
		newBalance := balance - amount
		a.Accounts[signer] = newBalance

		return &Tx{
			Action: Withdraw,
			Signer: signer,
			Amount: amount,
		}, nil
	}
}

func (a *Accounts) Send(sender string, recipient string, amount float64) ([]*Tx, error) {
	wtx, err := a.Withdraw(sender, amount)
	if err != nil {
		return nil, err
	}

	dtx := a.Deposit(recipient, amount)

	return []*Tx{wtx, dtx}, nil
}
