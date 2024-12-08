package server

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"sparrow/internal/block"
)

// Blocks implements block.Producer interface
type Blocks struct {
	logger  *zap.Logger
	config  *Config
	channel chan *block.Block

	stop chan struct{}
	done *sync.WaitGroup
}

var _ block.Producer = (*Blocks)(nil)

// NewBlocks
func NewBlocks(logger *zap.Logger, config *Config) (*Blocks, error) {
	return &Blocks{
		logger:  logger,
		config:  config,
		channel: make(chan *block.Block),

		stop: make(chan struct{}),
		done: &sync.WaitGroup{},
	}, nil
}

// Start implements block.Producer interface
func (bp *Blocks) Start(context.Context) error {
	bp.done.Add(1)

	go func() {
		defer bp.done.Done()

		var ticker = time.NewTicker(bp.config.Server.Ticker)

		for {
			select {
			case <-bp.stop:
				return
			case ts := <-ticker.C:
				bp.Produce(ts)
			}
		}
	}()

	return nil
}

// Stop implements block.Producer interface
func (bp *Blocks) Stop(context.Context) error {
	// send stop signal
	close(bp.stop)

	// wait done
	bp.done.Wait()

	// exit
	return nil
}

// Channel implements block.Producer interface
func (bp *Blocks) Channel() block.ProducerChannel {
	return bp.channel
}

// Produce
func (bp *Blocks) Produce(ts time.Time) {
	var block = block.NewBlock()

	defer bp.logger.Info("tick", zap.Time("time", ts), zap.Any("block", block))

	bp.channel <- block
}
