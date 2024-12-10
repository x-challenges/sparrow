package prices

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"sparrow/internal/instruments"
	"sparrow/internal/jupyter"
)

// Serice
type Service interface {
	// Update
	Update(ctx context.Context)
}

// Service interface implementation
type service struct {
	logger      *zap.Logger
	jupyter     jupyter.Client
	instruments instruments.Service
	rates       *Rates
}

// NewService
func newService(
	logger *zap.Logger,
	jupyter jupyter.Client,
	instruments instruments.Service,
	rates *Rates,
) (*service, error) {
	return &service{
		logger:      logger,
		jupyter:     jupyter,
		instruments: instruments,
		rates:       rates,
	}, nil
}

// Update implements Service interface
func (s *service) Update(ctx context.Context) {
	var (
		group, groupCtx = errgroup.WithContext(ctx)
		all             instruments.Instruments
		err             error
	)

	// take all instruments from storage
	if all, err = s.instruments.List(ctx); err != nil {
		s.logger.Error("take all instruments failed", zap.Error(err))
	}

	for _, base := range all {
		var base = *base

		group.Go(
			func() error {
				var (
					prices *jupyter.Prices
				)

				// take all prices
				if prices, err = s.jupyter.Prices(groupCtx, &base, all); err != nil {
					return err
				}

				// store prices
				for _, price := range prices.Data {
					s.rates.Store(base.Address, price.ID, price.Price)
				}

				return nil
			},
		)
	}

	// wait all goroutines
	if err = group.Wait(); err != nil {
		s.logger.Error("update price failed", zap.Error(err))
	}

	s.logger.Info("prices updated")
}
