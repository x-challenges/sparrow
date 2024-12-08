package server

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/x-challenges/raven/modules/config"

	"sparrow/internal/block"
	"sparrow/internal/instruments"
	"sparrow/internal/jupyter"
	"sparrow/internal/routes"
)

// ModuleName
const ModuleName = "server"

// Module
var Module = fx.Module(
	ModuleName,

	// internal
	jupyter.Module,
	instruments.Module,
	routes.Module,
	block.Module,

	config.Inject(new(Config)),

	// private usage
	fx.Provide(
		fx.Private,

		// server
		fx.Annotate(
			NewServer,
			fx.OnStart(func(ctx context.Context, s *Server) error { return s.Start(ctx) }),
			fx.OnStop(func(ctx context.Context, s *Server) error { return s.Stop(ctx) }),
		),

		// processor
		fx.Annotate(newProcessor, fx.As(new(Processor))),

		// block producer based on time ticker
		fx.Annotate(NewBlockProducer, fx.As(new(block.Producer))),
	),

	// force
	fx.Invoke(func(*Server) {}),

	fx.Decorate(
		func(logger *zap.Logger) *zap.Logger { return logger.Named(ModuleName) },
	),
)
