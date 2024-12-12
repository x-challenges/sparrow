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
	repository  Repository

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
	repository Repository,
) (*service, error) {
	return &service{
		logger:      logger,
		blocks:      blocks.Subscribe(),
		jupiter:     jupiter,
		instruments: instruments,
		config:      config,
		repository:  repository,

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

	go func() {
		var ctx = context.Background()

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
	}()

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
	for input := range s.instruments.Input(ctx) {
		for output := range s.instruments.Output(ctx) {
			s.repository.Store(ctx, input.Address, output.Address, 0) // init direct exchange rate
			s.repository.Store(ctx, output.Address, input.Address, 0) // init reverse exchange rate
		}
	}

	return nil
}

// Update
func (s *service) Update(ctx context.Context) {
	var (
		started         = time.Now()
		outputAddresses []string
		err             error
	)

	// take output addresses as a slice
	if outputAddresses, err = s.instruments.Addresses(ctx); err != nil {
		s.logger.Error("take output addresses failed", zap.Error(err))
		return
	}

	var (
		group, groupCtx = errgroup.WithContext(ctx)
		counter         = 0
	)

	for input := range s.instruments.Input(ctx) {
		var chunks = slices.Chunk(
			outputAddresses, s.config.Prices.Loader.ChunkSize)

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
						s.repository.Store(ctx, input.Address, price.ID, price.Price)   // direct rate
						s.repository.Store(ctx, price.ID, input.Address, 1/price.Price) // reverse rate

						counter += 2
					}

					return nil
				},
			)
		}
	}

	// wait all goroutines
	if err = group.Wait(); err != nil {
		s.logger.Error("update price failed", zap.Error(err))
		return
	}

	s.logger.Info("prices updated",
		zap.Duration("elapsed", time.Since(started)),
		zap.Int("count", counter),
	)
}

// Exchange implements Service interface
func (s *service) Exchange(_ context.Context, _, _ string, _ int64) (int64, error) {
	return 0, nil
}
