package quotes

import (
	"context"

	"go.uber.org/zap"

	"sparrow/internal/instruments"
	"sparrow/internal/jupyter"
)

// Service
type Service interface {
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
}

var _ Service = (*service)(nil)

// NewService
func newService(logger *zap.Logger, client jupyter.Client, repository Repository) (*service, error) {
	return &service{
		logger:     logger,
		client:     client,
		repository: repository,
	}, nil
}

// Quote implements Service interface
func (s *service) Quote(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quote, error) {
	return s.Quote(ctx, input, output, amount)
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

	// take reverse quotes
	amount = quotes.Direct.OutAmount

	if quotes.Reverse, err = s.Quote(ctx, output, input, amount); err != nil {
		return nil, err
	}

	return quotes, nil
}
