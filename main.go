package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kr/pretty"
)

func main() {
	r := chi.NewRouter()
	p := newTradingPlatform()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/", sayHello)

	// GET /orderbooks?base=ETH&qoute=USD
	r.Route("/orderbooks", func(chi.Router) {
		r.Get("/", p.handleGetOrderbooks)
	})
	// POST /orders
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", p.handleCreateOrder)
	})

	// GET accounts/:signer/balance
	// POST accounts/:signer/withdraw
	// POST accounts/:signer/deposit
	r.Route("/accounts", func(r chi.Router) {
		r.Route("/{signer}", func(r chi.Router) {
			r.Use(AccountCtx)
			r.Get("/", p.handleGetAccountBalance)
			r.Post("/", p.handleCreateAccount)
			r.Post("/{signer}/deposit", p.handleAccountDeposit)
			r.Post("/{signer}/withdraw", p.handleAccountWithdraw)
		})
	})

	pretty.Log("Starting server...")
	http.ListenAndServe(":8080", r)
}

func AccountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signer := chi.URLParam(r, "signer")
		ctx := context.WithValue(r.Context(), "signer", signer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
