package server

import "go.uber.org/zap"

// Processor
type Processor interface{}

// Processor interface implementation
type processor struct {
	logger *zap.Logger
}

var _ Processor = (*processor)(nil)

// NewProcessor
func newProcessor(logger *zap.Logger) *processor {
	return &processor{
		logger: logger,
	}
}
