package instruments

import (
	"context"
	"math"
	"math/big"
	"slices"

	"go.uber.org/zap"

	"sparrow/internal/jupiter"
)

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

	// Addresses
	Addresses(ctx context.Context) ([]string, error)
}

// Service interface implementation
type service struct {
	logger     *zap.Logger
	jupiter    jupiter.Client
	config     *Config
	repository Repository
}

var _ Service = (*service)(nil)

// NewService
func newService(logger *zap.Logger, jupiter jupiter.Client, config *Config, repository Repository) (*service, error) {
	return &service{
		logger:     logger,
		jupiter:    jupiter,
		config:     config,
		repository: repository,
	}, nil
}

// Start
func (s *service) Start(ctx context.Context) error {
	var (
		tokens jupiter.Tokens
		err    error
	)

	// try to load all tokens from jupiter
	if tokens, err = s.jupiter.Tokens(ctx); err != nil {
		return err
	}

	var (
		dailyVolume = s.config.Instruments.Loader.Skip.DailyVolume
		counter     int
	)

	// load all tokens as a instruments to storage
	for _, token := range tokens {
		// skip token by daily_volume
		if dailyVolume != 0 && token.DailyVolume <= dailyVolume {
			continue
		}

		var (
			instrument = &Instrument{
				Address:    token.Address,
				Ticker:     token.Name,
				Decimals:   token.Decimals,
				Zeros:      int64(math.Pow10(token.Decimals)),
				zerosValue: new(big.Float).SetInt64(int64(math.Pow10(token.Decimals))),
				token:      &token,
			}
		)

		// save instrument to storage
		if err = s.repository.Store(ctx, instrument); err != nil {
			return err
		}

		counter++
	}

	s.logger.Info("loaded tokens to instruments storage",
		zap.Int("total", len(tokens)),
		zap.Int("actual", counter),
	)

	// setup instruments from configuration
	for _, i := range s.config.Instruments.Pool {
		var (
			instrument *Instrument
			err        error
		)

		// load instrument for configuration
		if instrument, err = s.repository.Get(ctx, i.Address); err != nil {
			return err
		}

		instrument.Tags = i.Tags

		// store instrument into storage
		if err = s.repository.Store(ctx, instrument); err != nil {
			return err
		}
	}

	s.logger.Info("completed")

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

// Addresses
func (s *service) Addresses(ctx context.Context) ([]string, error) {
	var (
		res = make([]string, 0, 1000)
		all Instruments
		err error
	)

	if all, err = s.repository.List(ctx); err != nil {
		return nil, err
	}

	for _, inst := range all {
		res = append(res, inst.Address)
	}

	return res, nil
}
