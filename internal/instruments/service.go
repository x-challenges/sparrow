package instruments

import (
	"context"

	"go.uber.org/zap"
)

// Service
type Service interface {
	// Get
	Get(ctx context.Context, address string) (*Instrument, error)

	// List
	List(ctx context.Context) (Instruments, error)
}

// Service interface implementation
type service struct {
	logger     *zap.Logger
	repository Repository
}

var _ Service = (*service)(nil)

// NewService
func newService(logger *zap.Logger, repository Repository) (*service, error) {
	return &service{
		logger:     logger,
		repository: repository,
	}, nil
}

// Get implements Service interface
func (s *service) Get(ctx context.Context, address string) (*Instrument, error) {
	return s.repository.Get(ctx, address)
}

// List implements Service interface
func (s *service) List(ctx context.Context) (Instruments, error) {
	return s.repository.List(ctx)
}
