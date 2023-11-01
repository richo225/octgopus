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
	accounts.deposit("alice", 100)

	balance, err := accounts.balanceOf("alice")
	assert.NoError(t, err, "balanceOf(alice) should not return an error")
	assert.Equal(t, uint64(100), balance, "balanceOf(alice) should return 100")
}

func TestDepositNewAccount(t *testing.T) {
	accounts := newAccounts()

	// Test deposit with an empty account.
	err := accounts.deposit("alice", 100)
	assert.NoError(t, err, "deposit(alice, 100) should not return an error")

	balance, err := accounts.balanceOf("alice")
	assert.NoError(t, err, "balanceOf(alice) should not return an error")
	assert.Equal(t, uint64(100), balance, "balanceOf(alice) should return 100")
}

func TestDepositExistingAccount(t *testing.T) {
	accounts := newAccounts()

	// Test deposit with an existing non-zero balance.
	err := accounts.deposit("alice", 100)
	assert.NoError(t, err, "deposit(alice, 100) should not return an error")
	err = accounts.deposit("alice", 50)
	assert.NoError(t, err, "deposit(alice, 50) should not return an error")

	balance, err := accounts.balanceOf("alice")
	assert.NoError(t, err, "balanceOf(alice) should not return an error")
	assert.Equal(t, uint64(150), balance, "balanceOf(alice) should return 150")
}

func TestWithdrawAccountNotFound(t *testing.T) {
	accounts := newAccounts()

	err := accounts.withdraw("alice", 50)
	assert.Error(t, err, "withdraw(alice, 50) should return an error")
	assert.IsType(t, &AccountNotFoundError{}, err, "withdraw(alice, 50) should return an AccountNotFoundError")
}

func TestWithdrawInsufficientFunds(t *testing.T) {
	accounts := newAccounts()
	accounts.deposit("alice", 100)

	err := accounts.withdraw("alice", 150)
	assert.Error(t, err, "withdraw(alice, 150) should return an error")
	assert.IsType(t, &AccountUnderFundedError{}, err, "withdraw(alice, 150) should return an InsufficientFundsError")
}

func TestWithdrawSufficientFunds(t *testing.T) {
	accounts := newAccounts()
	accounts.deposit("alice", 100)

	err := accounts.withdraw("alice", 50)
	assert.NoError(t, err, "withdraw(alice, 50) should not return an error")

	balance, err := accounts.balanceOf("alice")
	assert.NoError(t, err, "balanceOf(alice) should not return an error")
	assert.Equal(t, uint64(50), balance, "balanceOf(alice) should return 50")
}

func TestSendWithSenderNotFound(t *testing.T) {
	accounts := newAccounts()

	err := accounts.send("alice", "bob", 50)
	assert.IsType(t, &AccountNotFoundError{}, err, "send(alice, bob, 50) should return an AccountNotFoundError")
}

func TestSendWithSenderUnderFunded(t *testing.T) {
	accounts := newAccounts()
	accounts.deposit("alice", 25)

	err := accounts.send("alice", "bob", 50)
	assert.IsType(t, &AccountUnderFundedError{}, err, "send(alice, bob, 50) should return an AccountUnderFundedError")
}

func TestSendWSuccess(t *testing.T) {
	accounts := newAccounts()
	accounts.deposit("alice", 100)
	accounts.deposit("bob", 10)

	err := accounts.send("alice", "bob", 30)
	assert.NoError(t, err, "send(alice, bob, 30) should not return an error")

	balance, _ := accounts.balanceOf("alice")
	assert.Equal(t, uint64(70), balance, "balanceOf(alice) should return 70")

	balance, _ = accounts.balanceOf("bob")
	assert.Equal(t, uint64(40), balance, "balanceOf(bob) should return 40")
}
