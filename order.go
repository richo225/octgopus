package main

import "time"

type Side string

const (
	Bid Side = "bid"
	Ask Side = "ask"
)

type Order struct {
	Side      Side   `json:"side"`
	Price     uint64 `json:"price"`
	Size      uint64 `json:"size"`
	Timestamp int64  `json:"timestamp"`
}

func newOrder(side Side, size uint64) *Order {
	return &Order{
		Side:      side,
		Size:      size,
		Timestamp: time.Now().UnixNano(),
	}
}
