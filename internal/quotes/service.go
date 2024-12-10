package quotes

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/x-challenges/raven/kun/model"

	"sparrow/internal/instruments"
	"sparrow/internal/jupyter"
)

// Service
type Service interface {
	// BatchInsert
	BatchInsert(ctx context.Context, quotes ...*Quotes) error

	// Quotes
	Quotes(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quotes, error)
}

// Service interface implementation
type service struct {
	logger     *zap.Logger
	client     jupyter.Client
	repository Repository
	config     *Config
}

var _ Service = (*service)(nil)

// NewService
func newService(logger *zap.Logger, client jupyter.Client, repository Repository, config *Config) (*service, error) {
	return &service{
		logger:     logger,
		client:     client,
		repository: repository,
		config:     config,
	}, nil
}

// Quotes implements Service interface
func (s *service) Quotes(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quotes, error) {
	var (
		group, groupCtx = errgroup.WithContext(ctx)
		quotes          = &Quotes{Direct: new(Quote), Reverse: new(Quote)}
		err             error
	)

	// take direct quote
	group.Go(
		func() error {
			quotes.Direct.StartedAt = time.Now().UnixMilli()
			quotes.Direct.Quote, err = s.client.Quote(groupCtx, input, output, amount, jupyter.WithSwapMode(jupyter.ExactIn))
			quotes.Direct.EndedAt = time.Now().UnixMilli()

			return err
		},
	)

	// take reverse quote
	group.Go(
		func() error {
			quotes.Reverse.StartedAt = time.Now().UnixMilli()
			quotes.Reverse.Quote, err = s.client.Quote(groupCtx, output, input, amount, jupyter.WithSwapMode(jupyter.ExactOut))
			quotes.Reverse.EndedAt = time.Now().UnixMilli()
			return err
		},
	)

	// wait group
	if err = group.Wait(); err != nil {
		return nil, err
	}

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
