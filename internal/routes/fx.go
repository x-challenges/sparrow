package routes

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/x-challenges/raven/modules/config"
)

var (
	start = func(ctx context.Context, s Service) error { return s.(*service).Start(ctx) }
	stop  = func(ctx context.Context, s Service) error { return s.(*service).Stop(ctx) }
)

// ModuleName
const ModuleName = "routes"

// Module
var Module = fx.Module(
	ModuleName,

	config.Inject(new(Config)),

	// public usage
	fx.Provide(
		fx.Annotate(
			newService,
			fx.As(new(Service)),
			fx.OnStart(start),
			fx.OnStop(stop),
		),
	),

	fx.Decorate(
		func(logger *zap.Logger) *zap.Logger { return logger.Named(ModuleName) },
	),
)
