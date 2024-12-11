package tokens

import (
	"context"

	"go.uber.org/zap"

	"sparrow/internal/jupyter"
)

// Service
type Service interface {
	// Get
	Get(ctx context.Context, address string) (*Token, error)
}

// Service interface implementation
type service struct {
	logger     *zap.Logger
	jupyter    jupyter.Client
	repository Repository
}

var _ Service = (*service)(nil)

// NewService
func newService(logger *zap.Logger, jupyter jupyter.Client, repository Repository) (*service, error) {
	return &service{
		logger:     logger,
		jupyter:    jupyter,
		repository: repository,
	}, nil
}

// Start
func (s *service) Start(ctx context.Context) error {
	var (
		tokens Tokens
		err    error
	)

	// try to load all tokens from jupyter
	if tokens, err = s.jupyter.Tokens(ctx); err != nil {
		return err
	}

	// save all tokens to storage
	for _, token := range tokens {
		if err = s.repository.Store(ctx, &token); err != nil {
			return err
		}
	}

	s.logger.Info("loaded all tokens info", zap.Int("count", len(tokens)))

	return nil
}

// Get implements Service interface
func (s *service) Get(ctx context.Context, address string) (*Token, error) {
	return s.repository.Get(ctx, address)
}
