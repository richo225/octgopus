package main

import (
	"fmt"
	"net/http"

	"github.com/kr/pretty"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	pretty.Print("Hello")
}

// Orderbooks
func (platform *TradingPlatform) handleGetOrderbooks(w http.ResponseWriter, r *http.Request) {

}

// Orders
func (platform *TradingPlatform) handleCreateOrder(w http.ResponseWriter, r *http.Request) {

}

// Accounting
func (platform *TradingPlatform) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	signer := r.Context().Value("signer").(string)
	err := platform.accounts.createAccount(signer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Write([]byte(fmt.Sprintf("Account created for %s", signer)))
}

func (platform *TradingPlatform) handleGetAccountBalance(w http.ResponseWriter, r *http.Request) {
	signer := r.Context().Value("signer").(string)
	balance, err := platform.accounts.balanceOf(signer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write([]byte(fmt.Sprintf("Balance: %d", balance)))
}

func (platform *TradingPlatform) handleAccountDeposit(w http.ResponseWriter, r *http.Request) {

}

func (platform *TradingPlatform) handleAccountWithdraw(w http.ResponseWriter, r *http.Request) {

}
