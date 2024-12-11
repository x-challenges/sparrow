package prices

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"sparrow/internal/block"
	"sparrow/internal/instruments"
	"sparrow/internal/jupyter"
)

// Service
type Service interface {
	// Update
	Update(ctx context.Context)
}

// Service interface implementation
type service struct {
	logger      *zap.Logger
	blocks      block.Listener
	jupyter     jupyter.Client
	instruments instruments.Service
	rates       *Rates

	stop chan struct{}
	done *sync.WaitGroup
}

// NewService
func newService(
	logger *zap.Logger,
	blocks block.Service,
	jupyter jupyter.Client,
	instruments instruments.Service,
	rates *Rates,
) (*service, error) {
	return &service{
		logger:      logger,
		blocks:      blocks.Subscribe(),
		jupyter:     jupyter,
		instruments: instruments,
		rates:       rates,

		stop: make(chan struct{}),
		done: &sync.WaitGroup{},
	}, nil
}

// Start
func (s *service) Start(ctx context.Context) error {
	s.done.Add(1)

	go func() {
		defer s.blocks.Close()
		defer s.done.Done()

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
	// send stop channel
	close(s.stop)

	// wait done
	s.done.Wait()

	// exit
	return nil
}

// Update implements Service interface
func (s *service) Update(context.Context) {
	// var (
	// 	group, groupCtx = errgroup.WithContext(ctx)
	// 	all             instruments.Iterator
	// 	err             error
	// )

	// take all instruments iterator
	// if all, err = s.instruments.All(ctx); err != nil {
	// 	s.logger.Error("take all instruments failed", zap.Error(err))
	// }

	// for base := range all {
	// 	var base = *base

	// 	group.Go(
	// 		func() error {
	// 			var (
	// 				prices *jupyter.Prices
	// 			)

	// 			// take all prices
	// 			if prices, err = s.jupyter.Prices(groupCtx, &base, all); err != nil {
	// 				return err
	// 			}

	// 			// store prices
	// 			for _, price := range prices.Data {
	// 				s.rates.Store(base.Address, price.ID, price.Price)
	// 			}

	// 			return nil
	// 		},
	// 	)
	// }

	// wait all goroutines
	// if err = group.Wait(); err != nil {
	// 	s.logger.Error("update price failed", zap.Error(err))
	// }

	s.logger.Info("prices updated")
}
