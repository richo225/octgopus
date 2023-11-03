package main

import "time"

type Side string

const (
	Bid Side = "bid"
	Ask Side = "ask"
)

type Order struct {
	side      Side
	limit     *Limit
	size      uint64
	timestamp int64
}

func newOrder(side Side, size uint64) *Order {
	return &Order{
		side:      side,
		size:      size,
		timestamp: time.Now().UnixNano(),
	}
}
