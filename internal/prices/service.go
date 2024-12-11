package prices

import (
	"context"
	"slices"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"sparrow/internal/block"
	"sparrow/internal/instruments"
	"sparrow/internal/jupiter"
)

// Service
type Service interface {
	// Exchange
	Exchange(ctx context.Context, input, output string, amount int64) (int64, error)
}

// Service interface implementation
type service struct {
	logger      *zap.Logger
	blocks      block.Listener
	jupiter     jupiter.Client
	instruments instruments.Service
	config      *Config
	rates       *Rates

	stop chan struct{}
	done *sync.WaitGroup
}

// NewService
func newService(
	logger *zap.Logger,
	blocks block.Service,
	jupiter jupiter.Client,
	instruments instruments.Service,
	config *Config,
	rates *Rates,
) (*service, error) {
	return &service{
		logger:      logger,
		blocks:      blocks.Subscribe(),
		jupiter:     jupiter,
		instruments: instruments,
		config:      config,
		rates:       rates,

		stop: make(chan struct{}),
		done: &sync.WaitGroup{},
	}, nil
}

// Start
func (s *service) Start(ctx context.Context) error {
	// init rates storage
	if err := s.Init(ctx); err != nil {
		return err
	}

	// run background tasks for updates
	s.done.Add(1)

	go func(ctx context.Context) {
		defer s.blocks.Close()
		defer s.done.Done()

		// first launch
		go s.Update(ctx)

		for {
			select {
			case <-s.stop:
				return
			case <-s.blocks.Updates():
				go s.Update(ctx)
			}
		}
	}(ctx)

	return nil
}

// Stop
func (s *service) Stop(context.Context) error {
	// send stop signal
	close(s.stop)

	// wait done
	s.done.Wait()

	// exit
	return nil
}

// Init
func (s *service) Init(ctx context.Context) error {
	var (
		inputIterator  instruments.Iterator
		outputIterator instruments.Iterator
		err            error
	)

	// take input instruments iterator
	if inputIterator, err = s.instruments.Base(ctx); err != nil {
		return err
	}

	// take output instruments iterator
	if outputIterator, err = s.instruments.All(ctx); err != nil {
		return err
	}

	// init storage
	for input := range inputIterator {
		for output := range outputIterator {
			s.rates.Store(input.Address, output.Address, 0)
		}
	}

	return nil
}

// Update
func (s *service) Update(ctx context.Context) {
	var (
		started         = time.Now()
		inputIterator   instruments.Iterator
		outputAddresses []string
		err             error
	)

	// take input instruments iterator
	if inputIterator, err = s.instruments.Base(ctx); err != nil {
		s.logger.Error("take input instruments iterator failed", zap.Error(err))
		return
	}

	// take output addresses as a slice
	if outputAddresses, err = s.instruments.Addresses(ctx); err != nil {
		s.logger.Error("take output addresses failed", zap.Error(err))
		return
	}

	var group, groupCtx = errgroup.WithContext(ctx)

	for input := range inputIterator {
		var chunks = slices.Chunk(outputAddresses, s.config.Prices.Loader.ChunkSize)

		for chunk := range chunks {
			group.Go(
				func() error {
					var prices *jupiter.Prices

					// take all prices from jupiter
					if prices, err = s.jupiter.Prices(groupCtx, input.Address, chunk...); err != nil {
						return err
					}

					// store prices
					for _, price := range prices.Data {
						s.rates.Store(input.Address, price.ID, price.Price)
					}

					return nil
				},
			)
		}
	}

	// wait all goroutines
	if err = group.Wait(); err != nil {
		s.logger.Error("update price failed", zap.Error(err))
	}

	s.logger.Info("price updates", zap.Duration("elapsed", time.Since(started)))
}

// Exchange implements Service interface
func (s *service) Exchange(_ context.Context, _, _ string, _ int64) (int64, error) {
	return 0, nil
}
