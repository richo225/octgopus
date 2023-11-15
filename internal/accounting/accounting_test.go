package accounting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceOfAccountNotFound(t *testing.T) {
	accounts := NewAccounts()

	_, err := accounts.BalanceOf("alice")
	assert.IsType(t, &AccountNotFoundError{}, err, "balanceOf(alice) should return an AccountNotFoundError")
}

func TestBalanceOfAccountExists(t *testing.T) {
	accounts := NewAccounts()
	accounts.CreateAccount("alice")
	accounts.Deposit("alice", 100)

	balance, err := accounts.BalanceOf("alice")
	assert.NoError(t, err, "balanceOf(alice) should not return an error")
	assert.Equal(t, float64(100), balance, "balanceOf(alice) should return 100")
}

func TestDepositAccountNotFound(t *testing.T) {
	accounts := NewAccounts()

	tx := accounts.Deposit("alice", 100)
	assert.Equal(t, Deposit, tx.Action, "deposit(alice) should return a Deposit Tx")
	assert.Equal(t, "alice", tx.Signer, "deposit(alice) should return a Tx with alice as signer")
	assert.Equal(t, float64(100), tx.Amount, "deposit(alice) should return a Tx with 100 as amount")

	balance, _ := accounts.BalanceOf("alice")
	assert.Equal(t, float64(100), balance, "balanceOf(alice) should return 100")

}

func TestDepositExistingAccount(t *testing.T) {
	accounts := NewAccounts()
	accounts.Accounts["alice"] = 50

	tx := accounts.Deposit("alice", 100)
	assert.Equal(t, Deposit, tx.Action, "deposit(alice) should return a Deposit Tx")
	assert.Equal(t, "alice", tx.Signer, "deposit(alice) should return a Tx with alice as signer")
	assert.Equal(t, float64(100), tx.Amount, "deposit(alice) should return a Tx with 100 as amount")

	balance, _ := accounts.BalanceOf("alice")
	assert.Equal(t, float64(150), balance, "balanceOf(alice) should return 150")
}

func TestWithdrawAccountNotFound(t *testing.T) {
	accounts := NewAccounts()

	_, err := accounts.Withdraw("alice", 50)
	assert.Error(t, err, "withdraw(alice, 50) should return an error")
	assert.IsType(t, &AccountNotFoundError{}, err, "withdraw(alice, 50) should return an AccountNotFoundError")
}

func TestWithdrawInsufficientFunds(t *testing.T) {
	accounts := NewAccounts()
	accounts.CreateAccount("alice")

	_, err := accounts.Withdraw("alice", 150)
	assert.Error(t, err, "withdraw(alice, 150) should return an error")
	assert.IsType(t, &AccountUnderFundedError{}, err, "withdraw(alice, 150) should return an InsufficientFundsError")
}

func TestWithdrawSufficientFunds(t *testing.T) {
	accounts := NewAccounts()
	accounts.CreateAccount("alice")
	accounts.Deposit("alice", 100)

	tx, err := accounts.Withdraw("alice", 50)
	assert.NoError(t, err, "withdraw(alice, 50) should not return an error")
	assert.Equal(t, Withdraw, tx.Action, "withdraw(alice, 50) should return a Withdraw Tx")
	assert.Equal(t, "alice", tx.Signer, "withdraw(alice, 50) should return a Tx with alice as signer")
	assert.Equal(t, float64(50), tx.Amount, "withdraw(alice, 50) should return a Tx with 50 as amount")

	balance, err := accounts.BalanceOf("alice")
	assert.NoError(t, err, "balanceOf(alice) should not return an error")
	assert.Equal(t, float64(50), balance, "balanceOf(alice) should return 50")
}

func TestSendWithSenderNotFound(t *testing.T) {
	accounts := NewAccounts()
	accounts.CreateAccount("bob")

	_, err := accounts.Send("alice", "bob", 50)
	assert.IsType(t, &AccountNotFoundError{}, err, "send(alice, bob, 50) should return an AccountNotFoundError")
}

func TestSendWithSenderUnderFunded(t *testing.T) {
	accounts := NewAccounts()
	accounts.CreateAccount("alice")
	accounts.Deposit("alice", 25)

	_, err := accounts.Send("alice", "bob", 50)
	assert.IsType(t, &AccountUnderFundedError{}, err, "send(alice, bob, 50) should return an AccountUnderFundedError")
}

func TestSendWSuccess(t *testing.T) {
	accounts := NewAccounts()
	accounts.CreateAccount("alice")
	accounts.CreateAccount("bob")
	accounts.Deposit("alice", 100)
	accounts.Deposit("bob", 10)

	tx, err := accounts.Send("alice", "bob", 30)
	assert.NoError(t, err, "send(alice, bob, 30) should not return an error")
	assert.Equal(t, Withdraw, tx[0].Action, "send(alice, bob, 30) should return a Withdraw Tx")
	assert.Equal(t, "alice", tx[0].Signer, "send(alice, bob, 30) should return a Tx with alice as signer")
	assert.Equal(t, float64(30), tx[0].Amount, "send(alice, bob, 30) should return a Tx with 30 as amount")

	assert.Equal(t, Deposit, tx[1].Action, "send(alice, bob, 30) should return a Deposit Tx")
	assert.Equal(t, "bob", tx[1].Signer, "send(alice, bob, 30) should return a Tx with bob as signer")
	assert.Equal(t, float64(30), tx[1].Amount, "send(alice, bob, 30) should return a Tx with 30 as amount")

	balance, _ := accounts.BalanceOf("alice")
	assert.Equal(t, float64(70), balance, "balanceOf(alice) should return 70")

	balance, _ = accounts.BalanceOf("bob")
	assert.Equal(t, float64(40), balance, "balanceOf(bob) should return 40")
}
