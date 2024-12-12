package instruments

import (
	"context"
	"slices"

	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
)

// Repository
type Repository interface {
	// Get
	Get(ctx context.Context, address string) (*Instrument, error)

	// Store
	Store(ctx context.Context, instrument *Instrument) error

	// Range
	Range(ctx context.Context, tag Tag) Iterator

	// Count
	Count(ctx context.Context, tag Tag) int
}

// Repository interface implementation
type repository struct {
	logger *zap.Logger
	data   *xsync.MapOf[string, *Instrument]
}

var _ Repository = (*repository)(nil)

// NewRepository
func newRepository(logger *zap.Logger) (*repository, error) {
	return &repository{
		logger: logger,
		data: xsync.NewMapOf[string, *Instrument](
			xsync.WithPresize(5000),
		),
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

// Store implements Service interface
func (rp *repository) Store(_ context.Context, instrument *Instrument) error {
	rp.data.Store(instrument.Address, instrument)
	return nil
}

// Range implements Repository interface
func (rp *repository) Range(ctx context.Context, tag Tag) Iterator {
	return func(yield func(*Instrument) bool) {
		rp.data.Range(func(_ string, value *Instrument) bool {
			if tag != Unspecified && !slices.Contains(value.Tags, tag) {
				return true
			}

			if !yield(value) {
				return true
			}

			return true
		})
	}
}

// Count implements Repository interface
func (rp *repository) Count(_ context.Context, tag Tag) int {
	var counter int

	rp.data.Range(func(_ string, value *Instrument) bool {
		if tag == Unspecified || slices.Contains(value.Tags, tag) {
			counter++
		}
		return true
	})

	return counter
}
