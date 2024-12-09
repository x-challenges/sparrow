package quotes

import "go.uber.org/zap"

// Repository
type Repository interface{}

// Repository interface implementation
type repository struct {
	logger *zap.Logger
}

var _ Repository = (*repository)(nil)

// NewRepository
func newRepository(logger *zap.Logger) (*repository, error) {
	return &repository{
		logger: logger,
	}, nil
}
