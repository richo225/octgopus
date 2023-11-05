package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceOfAccountNotFound(t *testing.T) {
	accounts := newAccounts()

	_, err := accounts.balanceOf("alice")
	assert.IsType(t, &AccountNotFoundError{}, err, "balanceOf(alice) should return an AccountNotFoundError")
}

func TestBalanceOfAccountExists(t *testing.T) {
	accounts := newAccounts()
	accounts.createAccount("alice")
	accounts.deposit("alice", 100)

	balance, err := accounts.balanceOf("alice")
	assert.NoError(t, err, "balanceOf(alice) should not return an error")
	assert.Equal(t, uint64(100), balance, "balanceOf(alice) should return 100")
}

func TestDepositAccountNotFound(t *testing.T) {
	accounts := newAccounts()

	_, err := accounts.deposit("alice", 100)
	assert.Error(t, err, "deposit(alice) should return an error")
	assert.IsType(t, &AccountNotFoundError{}, err, "deposit(alice) should return an AccountNotFoundError")
}

func TestDepositExistingAccount(t *testing.T) {
	accounts := newAccounts()
	accounts.createAccount("alice")

	tx, err := accounts.deposit("alice", 100)
	assert.NoError(t, err, "deposit(alice) should not return an error")
	assert.Equal(t, Deposit, tx.action, "deposit(alice) should return a Deposit Tx")
	assert.Equal(t, "alice", tx.signer, "deposit(alice) should return a Tx with alice as signer")
	assert.Equal(t, uint64(100), tx.amount, "deposit(alice) should return a Tx with 100 as amount")

	balance, _ := accounts.balanceOf("alice")
	assert.Equal(t, uint64(100), balance, "balanceOf(alice) should return 100")
}

func TestWithdrawAccountNotFound(t *testing.T) {
	accounts := newAccounts()

	_, err := accounts.withdraw("alice", 50)
	assert.Error(t, err, "withdraw(alice, 50) should return an error")
	assert.IsType(t, &AccountNotFoundError{}, err, "withdraw(alice, 50) should return an AccountNotFoundError")
}

func TestWithdrawInsufficientFunds(t *testing.T) {
	accounts := newAccounts()
	accounts.createAccount("alice")

	_, err := accounts.withdraw("alice", 150)
	assert.Error(t, err, "withdraw(alice, 150) should return an error")
	assert.IsType(t, &AccountUnderFundedError{}, err, "withdraw(alice, 150) should return an InsufficientFundsError")
}

func TestWithdrawSufficientFunds(t *testing.T) {
	accounts := newAccounts()
	accounts.createAccount("alice")
	accounts.deposit("alice", 100)

	tx, err := accounts.withdraw("alice", 50)
	assert.NoError(t, err, "withdraw(alice, 50) should not return an error")
	assert.Equal(t, Withdraw, tx.action, "withdraw(alice, 50) should return a Withdraw Tx")
	assert.Equal(t, "alice", tx.signer, "withdraw(alice, 50) should return a Tx with alice as signer")
	assert.Equal(t, uint64(50), tx.amount, "withdraw(alice, 50) should return a Tx with 50 as amount")

	balance, err := accounts.balanceOf("alice")
	assert.NoError(t, err, "balanceOf(alice) should not return an error")
	assert.Equal(t, uint64(50), balance, "balanceOf(alice) should return 50")
}

func TestSendWithSenderNotFound(t *testing.T) {
	accounts := newAccounts()
	accounts.createAccount("bob")

	_, err := accounts.send("alice", "bob", 50)
	assert.IsType(t, &AccountNotFoundError{}, err, "send(alice, bob, 50) should return an AccountNotFoundError")
}

func TestSendWithSenderUnderFunded(t *testing.T) {
	accounts := newAccounts()
	accounts.createAccount("alice")
	accounts.deposit("alice", 25)

	_, err := accounts.send("alice", "bob", 50)
	assert.IsType(t, &AccountUnderFundedError{}, err, "send(alice, bob, 50) should return an AccountUnderFundedError")
}

func TestSendWSuccess(t *testing.T) {
	accounts := newAccounts()
	accounts.createAccount("alice")
	accounts.createAccount("bob")
	accounts.deposit("alice", 100)
	accounts.deposit("bob", 10)

	tx, err := accounts.send("alice", "bob", 30)
	assert.NoError(t, err, "send(alice, bob, 30) should not return an error")
	assert.Equal(t, Withdraw, tx[0].action, "send(alice, bob, 30) should return a Withdraw Tx")
	assert.Equal(t, "alice", tx[0].signer, "send(alice, bob, 30) should return a Tx with alice as signer")
	assert.Equal(t, uint64(30), tx[0].amount, "send(alice, bob, 30) should return a Tx with 30 as amount")

	assert.Equal(t, Deposit, tx[1].action, "send(alice, bob, 30) should return a Deposit Tx")
	assert.Equal(t, "bob", tx[1].signer, "send(alice, bob, 30) should return a Tx with bob as signer")
	assert.Equal(t, uint64(30), tx[1].amount, "send(alice, bob, 30) should return a Tx with 30 as amount")

	balance, _ := accounts.balanceOf("alice")
	assert.Equal(t, uint64(70), balance, "balanceOf(alice) should return 70")

	balance, _ = accounts.balanceOf("bob")
	assert.Equal(t, uint64(40), balance, "balanceOf(bob) should return 40")
}
