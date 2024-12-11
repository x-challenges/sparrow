package instruments

import (
	"context"
	"iter"
	"slices"

	"go.uber.org/zap"
)

// Iterator
type Iterator iter.Seq[*Instrument]

// Service
type Service interface {
	// Get
	Get(ctx context.Context, address string) (*Instrument, error)

	// All
	All(ctx context.Context) (Iterator, error)

	// Base
	Base(ctx context.Context) (Iterator, error)

	// Routable
	Routable(ctx context.Context) (Iterator, error)

	// Swappable
	Swappable(ctx context.Context) (Iterator, error)
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

// Range
func (s *service) Range(ctx context.Context, tag Tag) (Iterator, error) {
	var (
		all Instruments
		err error
	)

	// take all instruments
	if all, err = s.repository.List(ctx); err != nil {
		return nil, err
	}

	return func(yield func(*Instrument) bool) {
		for _, instrument := range all {
			if tag != Unspecified && !slices.Contains(instrument.Tags, tag) {
				continue
			}

			if !yield(instrument) {
				return
			}
		}
	}, nil
}

// All implements Service interface
func (s *service) All(ctx context.Context) (Iterator, error) {
	return s.Range(ctx, Unspecified)
}

// Base implements Service interface
func (s *service) Base(ctx context.Context) (Iterator, error) {
	return s.Range(ctx, Base)
}

// Routable implements Service interface
func (s *service) Routable(ctx context.Context) (Iterator, error) {
	return s.Range(ctx, Route)
}

// Swappable implements Service interface
func (s *service) Swappable(ctx context.Context) (Iterator, error) {
	return s.Range(ctx, Swap)
}
