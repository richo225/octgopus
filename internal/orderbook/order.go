package orderbook

import (
	"encoding/json"
	"errors"
	"time"
)

type Side string

const (
	Bid Side = "bid"
	Ask Side = "ask"
)

func (side *Side) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case string(Bid), string(Ask):
		*side = Side(s)
	default:
		return errors.New("invalid side")
	}

	return nil
}

type Order struct {
	Side      Side    `json:"side"`
	Price     float64 `json:"price"`
	Size      float64 `json:"size"`
	Timestamp int64   `json:"timestamp"`
}

func NewOrder(side Side, size float64) *Order {
	return &Order{
		Side:      side,
		Size:      size,
		Timestamp: time.Now().UnixNano(),
	}
}
