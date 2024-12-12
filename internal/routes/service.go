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
	Range() Iterator
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

// Load
func (s *service) load(ctx context.Context) error {
	var (
		step     = s.config.Routes.Step
		from     = s.config.Routes.Range[0] * step
		to       = s.config.Routes.Range[1] * step
		priority = 0
	)

	// iterate base instruments
	for base := range s.instruments.Input(ctx) {

		// iterate quote instruments
		for quote := range s.instruments.Routable(ctx) {

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
					Amount:   new(big.Float).SetFloat64(float64(amount) / float64(step)),
					Priority: priority,
				}

				// insert route to the pool
				s.pool.AddRoute(route)
			}
		}
	}

	s.logger.Info("instruments loaded into route pool",
		zap.Int("count", len(s.pool.index)),
	)

	return nil
}

// Range implements Service interface
func (s *service) Range() Iterator {
	return s.pool.Range()
}
