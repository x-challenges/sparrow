package block

import (
	"time"

	"github.com/google/uuid"
)

// Status
type Status string

const (
	New  Status = "NEW"
	Done Status = "DONE"
)

// Block
type Block struct {
	ID        string    `json:"id"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// NewBlock
func NewBlock() *Block {
	return &Block{
		ID:        uuid.NewString(),
		Status:    New,
		CreatedAt: time.Now(),
	}
}
