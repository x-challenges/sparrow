package server

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"sparrow/internal/block"
)

// BlockProducer implements block.Producer interface
type BlockProducer struct {
	logger  *zap.Logger
	config  *Config
	channel chan *block.Block

	stop chan struct{}
	done *sync.WaitGroup
}

var _ block.Producer = (*BlockProducer)(nil)

// NewBlockProducer
func NewBlockProducer(logger *zap.Logger, config *Config) (*BlockProducer, error) {
	return &BlockProducer{
		logger:  logger,
		config:  config,
		channel: make(chan *block.Block),

		stop: make(chan struct{}),
		done: &sync.WaitGroup{},
	}, nil
}

// Start implements block.Producer interface
func (bp *BlockProducer) Start(context.Context) error {
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
func (bp *BlockProducer) Stop(context.Context) error {
	// send stop signal
	close(bp.stop)

	// wait done
	bp.done.Wait()

	// exit
	return nil
}

// Channel implements block.Producer interface
func (bp *BlockProducer) Channel() block.ProducerChannel { return bp.channel }

// Produce
func (bp *BlockProducer) Produce(ts time.Time) {
	var block = block.NewBlock()

	defer bp.logger.Info("tick", zap.Time("time", ts), zap.Any("block", block))

	bp.channel <- block
}
