package instruments

import (
	"context"
	"iter"
	"math"
	"math/big"
	"slices"

	"go.uber.org/zap"

	"sparrow/internal/tokens"
)

// Iterator
type Iterator iter.Seq[*Instrument]

// Service
type Service interface {
	// Get
	Get(ctx context.Context, address string) (*Instrument, error)

	// All
	All(ctx context.Context) (Iterator, error)

	// Base
	Base(ctx context.Context) (Iterator, error)

	// Routable
	Routable(ctx context.Context) (Iterator, error)
}

// Service interface implementation
type service struct {
	logger     *zap.Logger
	tokens     tokens.Service
	config     *Config
	repository Repository
}

var _ Service = (*service)(nil)

// NewService
func newService(logger *zap.Logger, tokens tokens.Service, config *Config, repository Repository) (*service, error) {
	return &service{
		logger:     logger,
		tokens:     tokens,
		config:     config,
		repository: repository,
	}, nil
}

// Start
func (s *service) Start(ctx context.Context) error {
	var err error

	// load all available instruments from config to inmemory map
	for _, instrument := range s.config.Instruments.Pool {
		var (
			zeros      = int64(math.Pow10(instrument.Decimals))
			instrument = &Instrument{
				Address:    instrument.Address,
				Ticker:     instrument.Ticker,
				Decimals:   instrument.Decimals,
				Tags:       instrument.Tags,
				Zeros:      zeros,
				zerosValue: new(big.Float).SetInt64(zeros),
			}
		)

		// store instrument into storage
		if err = s.repository.Store(ctx, instrument); err != nil {
			return err
		}
	}

	return nil
}

// Get implements Service interface
func (s *service) Get(ctx context.Context, address string) (*Instrument, error) {
	return s.repository.Get(ctx, address)
}

// Range
func (s *service) Range(ctx context.Context, tag Tag) (Iterator, error) {
	var (
		all Instruments
		err error
	)

	// take all instruments
	if all, err = s.repository.List(ctx); err != nil {
		return nil, err
	}

	return func(yield func(*Instrument) bool) {
		for _, instrument := range all {
			if tag != Unspecified && !slices.Contains(instrument.Tags, tag) {
				continue
			}

			if !yield(instrument) {
				return
			}
		}
	}, nil
}

// All implements Service interface
func (s *service) All(ctx context.Context) (Iterator, error) {
	return s.Range(ctx, Unspecified)
}

// Base implements Service interface
func (s *service) Base(ctx context.Context) (Iterator, error) {
	return s.Range(ctx, Base)
}

// Routable implements Service interface
func (s *service) Routable(ctx context.Context) (Iterator, error) {
	return s.Range(ctx, Route)
}
