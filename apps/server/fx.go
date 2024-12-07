package server

import (
	"context"

	"github.com/x-challenges/raven/modules/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"sparrow/internal/instruments"
	"sparrow/internal/jupyter"
	"sparrow/internal/routes"
)

var (
	// start server fn
	start = func(ctx context.Context, s *Server) error { return s.Start(ctx) }

	// stop server fn
	stop = func(ctx context.Context, s *Server) error { return s.Stop(ctx) }
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

	config.Inject(new(Config)),

	// private usage
	fx.Provide(
		fx.Private,

		fx.Annotate(
			NewServer,
			fx.OnStart(start),
			fx.OnStop(stop),
		),
	),

	// force
	fx.Invoke(func(*Server) {}),

	fx.Decorate(
		func(logger *zap.Logger) *zap.Logger { return logger.Named(ModuleName) },
	),
)
