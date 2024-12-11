package instruments

import (
	"context"

	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
)

// Repository
type Repository interface {
	// Get
	Get(ctx context.Context, address string) (*Instrument, error)

	// Store
	Store(ctx context.Context, instrument *Instrument) error

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
