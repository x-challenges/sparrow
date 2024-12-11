package tokens

import (
	"context"

	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
)

// Repository
type Repository interface {
	// Get
	Get(ctx context.Context, address string) (*Token, error)

	// Store
	Store(ctx context.Context, token *Token) error
}

// Repository interface implementation
type repository struct {
	logger *zap.Logger
	data   *xsync.MapOf[string, *Token]
}

var _ Repository = (*repository)(nil)

// NewRepository
func newRepository(logger *zap.Logger) (*repository, error) {
	return &repository{
		logger: logger,
		data: xsync.NewMapOf[string, *Token](
			xsync.WithPresize(5000),
		),
	}, nil
}

// Get implements Repository interface
func (rp *repository) Get(_ context.Context, address string) (*Token, error) {
	var (
		token *Token
		exist bool
	)

	// try to load token by address
	if token, exist = rp.data.Load(address); !exist {
		return nil, ErrNotFound
	}

	return token, nil
}

// Store implements Repository interface
func (rp *repository) Store(_ context.Context, token *Token) error {
	rp.data.Store(token.Address, token)
	return nil
}
