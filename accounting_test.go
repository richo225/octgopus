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
