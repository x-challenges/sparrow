package routes

import (
	"context"
	"math/big"

	"go.uber.org/zap"

	"sparrow/internal/instruments"
)

// Service
type Service interface {
	// Range
	Range(ctx context.Context) Iterator
}

// Service interface implementation
type service struct {
	logger      *zap.Logger
	instruments instruments.Service
	config      *Config
	pool        *Pool
}

var _ Service = (*service)(nil)

// NewService
func newService(logger *zap.Logger, instruments instruments.Service, config *Config) (*service, error) {
	return &service{
		logger:      logger,
		instruments: instruments,
		config:      config,
		pool:        NewPool(),
	}, nil
}

// Start
func (s *service) Start(ctx context.Context) error {
	return s.load(ctx)
}

// Stop
func (s *service) Stop(context.Context) error {
	return nil
}

// Load
func (s *service) load(ctx context.Context) error {
	var (
		allInstruments instruments.Instruments
		err            error
	)

	// take all available instruments
	if allInstruments, err = s.instruments.List(ctx); err != nil {
		return err
	}

	// all instruments
	var (
		step     = s.config.Routes.Step
		from     = s.config.Routes.Range[0] * step
		to       = s.config.Routes.Range[1] * step
		priority = 0
	)

	// iterate based instrument
	for _, baseAddr := range s.config.Routes.BaseCcy {
		var base *instruments.Instrument

		// take base instrument
		if base, err = s.instruments.Get(ctx, baseAddr); err != nil {
			return err
		}

		// iterate quoted instruments
		for _, quote := range allInstruments {

			// skip
			if base.Address == quote.Address {
				continue
			}

			// iterate all price range
			for amount := from; amount < to; amount++ {

				// skip if zero
				if amount == 0 {
					continue
				}

				var route = &Route{
					Base:     base,
					Quote:    quote,
					Amount:   new(big.Rat).SetFloat64(float64(amount) / float64(step)),
					Priority: priority,
				}

				s.pool.AddRoute(route)
			}
		}
	}

	s.logger.Info("instruments loaded into routes pool",
		zap.Int("count", len(s.pool.index)),
	)

	return nil
}

// Range implements Service interface
func (s *service) Range(context.Context) Iterator {
	return s.pool.Range()
}
