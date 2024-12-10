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
		data:   xsync.NewMapOf[string, float32](),
	}
}

// Key
func (r *Rates) Key(from, to string) string {
	return from + "/" + to
}

// Store
func (r *Rates) Store(from, to string, value float32) {
	r.data.Store(r.Key(from, to), value)
}

// Load
func (r *Rates) Load(from, to string) (float32, error) {
	if res, exist := r.data.Load(r.Key(from, to)); exist {
		return res, nil
	}

	return 0, ErrNotFound
}
