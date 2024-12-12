package quotes

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/x-challenges/raven/kun/model"

	"sparrow/internal/instruments"
	"sparrow/internal/jupiter"
	"sparrow/internal/prices"
)

// Service
type Service interface {
	// Quotes
	Quotes(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quotes, error)

	// BatchInsert
	BatchInsert(ctx context.Context, quotes ...*Quotes) error
}

// Service interface implementation
type service struct {
	logger     *zap.Logger
	jupiter    jupiter.Client
	prices     prices.Service
	config     *Config
	repository Repository
}

var _ Service = (*service)(nil)

// NewService
func newService(
	logger *zap.Logger,
	jupiter jupiter.Client,
	prices prices.Service,
	config *Config,
	repository Repository,
) (*service, error) {
	return &service{
		logger:     logger,
		jupiter:    jupiter,
		prices:     prices,
		config:     config,
		repository: repository,
	}, nil
}

// Quotes implements Service interface
func (s *service) Quotes(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quotes, error) {
	var (
		start           = time.Now()
		quotes          = new(Quotes)
		group, groupCtx = errgroup.WithContext(ctx)
		exchanged       float64
		err             error
	)

	// exchange
	if exchanged, err = s.prices.Exchange(groupCtx, input.Address, output.Address, float64(amount)); err != nil {
		return nil, err
	}

	// take direct quote
	group.Go(
		func() error {
			quotes.Direct, err = s.jupiter.Quote(groupCtx,
				input.Address, output.Address,
				int64(exchanged),
				jupiter.WithSwapMode(jupiter.ExactOut),
			)

			return err
		},
	)

	// take reverse quote
	group.Go(
		func() error {
			quotes.Reverse, err = s.jupiter.Quote(groupCtx,
				output.Address, input.Address,
				int64(exchanged),
				jupiter.WithSwapMode(jupiter.ExactIn),
			)

			return err
		},
	)

	if err = group.Wait(); err != nil {
		if errors.Is(err, jupiter.ErrRouteNotFound) {
			return nil, errors.Join(err, ErrQuoteNotFound)
		}
		return nil, err
	}

	// diff
	quotes.Diff = quotes.Reverse.OutAmount - quotes.Direct.InAmount

	// profit
	quotes.Profit = (1.0 - float32(quotes.Direct.InAmount)/float32(quotes.Reverse.OutAmount)) * 100.0

	// elapsed times
	quotes.Elapsed = time.Since(start)

	return quotes, nil
}

// BatchInsert implements Service interface
func (s *service) BatchInsert(ctx context.Context, quotes ...*Quotes) error {
	var instances = []*Model{}

	for _, vq := range quotes {
		// skip empty
		if vq == nil {
			continue
		}

		instances = append(instances,
			&Model{
				Base:   model.NewBase(),
				Quotes: *vq,
			},
		)
	}

	return s.repository.BatchInsert(ctx, instances...)
}
