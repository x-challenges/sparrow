package instruments

import (
	"context"
	"math"
	"math/big"

	"go.uber.org/zap"

	"sparrow/internal/jupiter"
)

// Service
type Service interface {
	// Get
	Get(ctx context.Context, address string) (*Instrument, error)

	// Input
	Input(ctx context.Context) Iterator

	// Output
	Output(ctx context.Context) Iterator

	// Routable
	Routable(ctx context.Context) Iterator

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

	s.logger.Info("instruments loaded",
		zap.Int("base", s.repository.Count(ctx, Base)),         // inputable
		zap.Int("route", s.repository.Count(ctx, Route)),       // outputable
		zap.Int("total", s.repository.Count(ctx, Unspecified)), // total
	)

	// metrics
	instrumentsCounter.Add(ctx, int64(s.repository.Count(ctx, Unspecified)))

	return nil
}

// Get implements Service interface
func (s *service) Get(ctx context.Context, address string) (*Instrument, error) {
	return s.repository.Get(ctx, address)
}

// Input implements Service interface
func (s *service) Input(ctx context.Context) Iterator {
	return s.repository.Range(ctx, Base)
}

// Output implements Service interface
func (s *service) Output(ctx context.Context) Iterator {
	return s.repository.Range(ctx, Unspecified)
}

// Routable implements Service interface
func (s *service) Routable(ctx context.Context) Iterator {
	return s.repository.Range(ctx, Route)
}

// Addresses
func (s *service) Addresses(ctx context.Context) ([]string, error) {
	var (
		res = make([]string, 0, 1000)
	)

	for inst := range s.repository.Range(ctx, Unspecified) {
		res = append(res, inst.Address)
	}

	return res, nil
}
