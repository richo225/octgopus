package main

import "time"

type Side string

const (
	Bid Side = "bid"
	Ask Side = "ask"
)

type Order struct {
	Side      Side    `json:"side"`
	Price     float64 `json:"price"`
	Size      float64 `json:"size"`
	Timestamp int64   `json:"timestamp"`
}

func newOrder(side Side, size float64) *Order {
	return &Order{
		Side:      side,
		Size:      size,
		Timestamp: time.Now().UnixNano(),
	}
}
