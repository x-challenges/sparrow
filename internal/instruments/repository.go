package instruments

import (
	"context"
	"math"
	"math/big"

	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
)

// Repository
type Repository interface {
	// Get
	Get(ctx context.Context, address string) (*Instrument, error)

	// List
	List(ctx context.Context) (Instruments, error)
}

// Repository interface implementation
type repository struct {
	logger *zap.Logger
	data   *xsync.MapOf[string, *Instrument]
}

var _ Repository = (*repository)(nil)

// NewRepository
func newRepository(logger *zap.Logger, config *Config) (*repository, error) {
	var data = xsync.NewMapOf[string, *Instrument](
		xsync.WithPresize(1000),
	)

	// load all available instruments from config to inmemory map
	for _, instrument := range config.Instruments.Pool {
		var zeros = int64(math.Pow10(instrument.Decimals))

		data.Store(
			// key
			instrument.Address,

			// data
			&Instrument{
				Address:    instrument.Address,
				Ticker:     instrument.Ticker,
				Decimals:   instrument.Decimals,
				Tags:       instrument.Tags,
				zeros:      zeros,
				zerosValue: new(big.Float).SetInt64(zeros),
			},
		)
	}

	return &repository{
		logger: logger,
		data:   data,
	}, nil
}

// Get implements Repository interface
func (rp *repository) Get(_ context.Context, address string) (*Instrument, error) {
	var (
		instrument *Instrument
		exists     bool
	)

	// try load
	if instrument, exists = rp.data.Load(address); !exists {
		return nil, ErrNotFound
	}

	return instrument, nil
}

// List implements Repository interface
func (rp *repository) List(_ context.Context) (Instruments, error) {
	var instruments = Instruments{}

	rp.data.Range(
		func(_ string, value *Instrument) bool {
			instruments = append(instruments, value)
			return true
		},
	)

	return instruments, nil
}
