package quotes

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	r "github.com/x-challenges/raven/kun/repository"
)

// Repository
type Repository interface {
	// BatchInsert
	BatchInsert(ctx context.Context, instances ...*Model) error
}

// Repository interface implementation
type repository struct {
	logger *zap.Logger
	db     *gorm.DB

	r.BatchInsertOp[*Model]
}

var _ Repository = (*repository)(nil)

// NewRepository
func newRepository(logger *zap.Logger, db *gorm.DB) (*repository, error) {
	var rdb = func(context.Context) *gorm.DB { return db }

	return &repository{
		logger: logger,
		db:     db,

		BatchInsertOp: r.NewBatchInsertOp[*Model](rdb),
	}, nil
}

// BatchInsert implements Repository interface
func (r *repository) BatchInsert(ctx context.Context, instances ...*Model) error {
	return r.BatchInsertOp.BatchInsert(ctx, instances...)
}
