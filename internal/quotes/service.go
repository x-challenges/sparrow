package quotes

import (
	"context"

	"go.uber.org/zap"

	"github.com/x-challenges/raven/kun/model"

	"sparrow/internal/instruments"
	"sparrow/internal/jupyter"
)

// Service
type Service interface {
	// BatchInsert
	BatchInsert(ctx context.Context, quotes ...*Quotes) error

	// Quote
	Quote(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quote, error)

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

// Quote implements Service interface
func (s *service) Quote(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quote, error) {
	return s.client.Quote(ctx, input, output, amount)
}

// Quotes implements Service interface
func (s *service) Quotes(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quotes, error) {
	var (
		quotes = new(Quotes)
		err    error
	)

	// take direct quotes
	if quotes.Direct, err = s.Quote(ctx, input, output, amount); err != nil {
		return nil, err
	}

	amount = quotes.Direct.OutAmount

	// take reverse quotes
	if quotes.Reverse, err = s.Quote(ctx, output, input, amount); err != nil {
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
