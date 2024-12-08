package block

import (
	"context"
)

// ProducerChannel
type ProducerChannel <-chan *Block

// Producer
type Producer interface {
	// Start
	Start(context.Context) error

	// Stop
	Stop(context.Context) error

	// Channel
	Channel() ProducerChannel
}
