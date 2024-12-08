package orders

import "go.uber.org/zap"

// Service
type Service interface{}

// Service interface implementation
type service struct {
	logger *zap.Logger
}

var _ Service = (*service)(nil)

// NewService
func newService(logger *zap.Logger) (*service, error) {
	return &service{
		logger: logger,
	}, nil
}
