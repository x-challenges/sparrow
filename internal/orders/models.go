package orders

import (
	"time"

	"github.com/google/uuid"
)

// Status
type Status = string

const (
	Pending Status = "PENDING"
	Success Status = "SUCCESS"
	Failed  Status = "FAILED"
)

// MetaData
type MetaData map[string]interface{}

// Order
type Order struct {
	ID         string    `json:"id"`
	Status     Status    `json:"status"`
	BaseCcy    string    `json:"base_ccy"`
	QuoteCcy   string    `json:"quote_ccy"`
	Price      int       `json:"price"`
	Amount     int       `json:"amount"`
	MetaData   MetaData  `json:"meta_data"`
	ExecutedAt time.Time `json:"executed_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewOrder
func NewOrder(baseCcy, quoteCcy string) *Order {
	return &Order{
		ID:        uuid.NewString(),
		Status:    Pending,
		BaseCcy:   baseCcy,
		QuoteCcy:  quoteCcy,
		CreatedAt: time.Now(),
	}
}
