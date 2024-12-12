package prices

import (
	"context"

	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
)

// Repository
type Repository interface {
	// Load
	Load(ctx context.Context, input, output string) (float64, error)

	// Store
	Store(ctx context.Context, input, output string, value float64)
}

// Repository interface implementation
type repository struct {
	logger *zap.Logger
	data   *xsync.MapOf[string, float64]
}

var _ Repository = (*repository)(nil)

// NewRepository
func newRepository(logger *zap.Logger) (*repository, error) {
	return &repository{
		logger: logger,
		data: xsync.NewMapOf[string, float64](
			xsync.WithPresize(5000),
			xsync.WithGrowOnly(),
		),
	}, nil
}

// Key
func (rp *repository) Key(input, output string) string {
	return input + "/" + output
}

// Store implements Repository interface
func (rp *repository) Store(_ context.Context, input, output string, value float64) {
	rp.data.Store(rp.Key(input, output), value)
}

// Load implements Repository interface
func (rp *repository) Load(_ context.Context, input, output string) (float64, error) {
	if res, exist := rp.data.Load(rp.Key(input, output)); exist {
		return res, nil
	}

	return 0, ErrNotFound
}
