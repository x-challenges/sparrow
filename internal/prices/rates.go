package prices

import (
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
)

// Rate
type Rates struct {
	logger *zap.Logger
	data   *xsync.MapOf[string, float32]
}

// Rates
func NewRates(logger *zap.Logger) *Rates {
	return &Rates{
		logger: logger,
		data: xsync.NewMapOf[string, float32](
			xsync.WithPresize(5000),
		),
	}
}

// Key
func (r *Rates) Key(input, output string) string {
	return input + "/" + output
}

// Store
func (r *Rates) Store(input, output string, value float32) {
	r.data.Store(r.Key(input, output), value)
}

// Load
func (r *Rates) Load(input, output string) (float32, error) {
	if res, exist := r.data.Load(r.Key(input, output)); exist {
		return res, nil
	}

	return 0, ErrNotFound
}
