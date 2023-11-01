package main

import (
	"testing"
)

func TestBalanceOfAccountNotFound(t *testing.T) {
	accounts := newAccounts()

	_, err := accounts.balanceOf("alice")
	if err == nil {
		t.Errorf("balanceOf(alice) returned nil error, expected AccountNotFoundError")
	} else if _, ok := err.(*AccountNotFoundError); !ok {
		t.Errorf("balanceOf(alice) returned %v, expected AccountNotFoundError", err)
	}
}

func TestBalanceOfAccountExists(t *testing.T) {
	accounts := newAccounts()
	accounts.deposit("alice", 100)

	balance, err := accounts.balanceOf("alice")
	if err != nil {
		t.Errorf("balanceOf(alice) returned error: %v", err)
	} else if balance != 100 {
		t.Errorf("balanceOf(alice) = %d, expected 100", balance)
	}
}

func TestDepositNewAccount(t *testing.T) {
	accounts := newAccounts()

	// Test deposit with an empty account.
	err := accounts.deposit("alice", 100)
	if err != nil {
		t.Errorf("deposit(alice, 100) returned error: %v", err)
	}
	balance, err := accounts.balanceOf("alice")
	if err != nil {
		t.Errorf("balanceOf(alice) returned error: %v", err)
	} else if balance != 100 {
		t.Errorf("balanceOf(alice) = %d, expected 100", balance)
	}
}

func TestDepositExistingAccount(t *testing.T) {
	accounts := newAccounts()

	// Test deposit with an existing non-zero balance.
	err := accounts.deposit("alice", 100)
	if err != nil {
		t.Errorf("deposit(alice, 100) returned error: %v", err)
	}
	err = accounts.deposit("alice", 50)
	if err != nil {
		t.Errorf("deposit(alice, 50) returned error: %v", err)
	}
	balance, err := accounts.balanceOf("alice")
	if err != nil {
		t.Errorf("balanceOf(alice) returned error: %v", err)
	} else if balance != 150 {
		t.Errorf("balanceOf(alice) = %d, expected 150", balance)
	}
}

func TestWithdrawAccountNotFound(t *testing.T) {
	accounts := newAccounts()

	err := accounts.withdraw("alice", 50)
	if err == nil {
		t.Errorf("withdraw(alice, 50) returned nil error, expected AccountNotFoundError")
	} else if _, ok := err.(*AccountNotFoundError); !ok {
		t.Errorf("withdraw(alice, 50) returned %v, expected AccountNotFoundError", err)
	}
}

func TestWithdrawUnderFunded(t *testing.T) {
	accounts := newAccounts()
	accounts.deposit("alice", 10)

	// Test withdraw with an empty account.
	err := accounts.withdraw("alice", 50)
	if err == nil {
		t.Errorf("withdraw(alice, 50) returned nil error, expected UnderFunded error")
	} else if _, ok := err.(*AccountUnderFundedError); !ok {
		t.Errorf("withdraw(alice, 50) returned %v, expected UnderFundedError", err)
	}
}

func TestWithdrawSufficientBalance(t *testing.T) {
	accounts := newAccounts()
	accounts.deposit("alice", 100)

	err := accounts.withdraw("alice", 20)
	if err != nil {
		t.Errorf("withdraw(alice, 50) returned error: %v", err)
	}
	balance, err := accounts.balanceOf("alice")
	if err != nil {
		t.Errorf("balanceOf(alice) returned error: %v", err)
	} else if balance != 80 {
		t.Errorf("balanceOf(alice) = %d, expected 80", balance)
	}
}

func TestSendWithSenderNotFound(t *testing.T) {
	accounts := newAccounts()

	err := accounts.send("alice", "bob", 50)
	if err == nil {
		t.Errorf("send(alice, bob, 50) returned nil error, expected AccountNotFoundError")
	} else if _, ok := err.(*AccountNotFoundError); !ok {
		t.Errorf("send(alice, bob, 50) returned %v, expected AccountNotFoundError", err)
	}
}

func TestSendWithSenderUnderFunded(t *testing.T) {
	accounts := newAccounts()
	accounts.deposit("alice", 25)

	err := accounts.send("alice", "bob", 50)
	if err == nil {
		t.Errorf("send(alice, bob, 50) returned nil error, expected AccountUnderFundedError")
	} else if _, ok := err.(*AccountUnderFundedError); !ok {
		t.Errorf("send(alice, bob, 50) returned %v, expected AccountUnderFundedError", err)
	}
}

func TestSendWSuccess(t *testing.T) {
	accounts := newAccounts()
	accounts.deposit("alice", 100)
	accounts.deposit("bob", 10)

	err := accounts.send("alice", "bob", 30)
	if err != nil {
		t.Errorf("send(alice, bob, 50) returned error: %v", err)
	}

	balance, err := accounts.balanceOf("alice")
	if err != nil {
		t.Errorf("balanceOf(alice) returned error: %v", err)
	} else if balance != 70 {
		t.Errorf("balanceOf(alice) = %d, expected 70", balance)
	}

	balance, err = accounts.balanceOf("bob")
	if err != nil {
		t.Errorf("balanceOf(bob) returned error: %v", err)
	} else if balance != 40 {
		t.Errorf("balanceOf(bob) = %d, expected 40", balance)
	}
}
