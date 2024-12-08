package block

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ModuleName
const ModuleName = "block"

var Module = fx.Module(
	ModuleName,

	// public usage
	fx.Provide(
		fx.Annotate(
			newService,
			fx.As(new(Service)),
			fx.OnStart(func(ctx context.Context, s Service) error { return s.(*service).Start(ctx) }),
			fx.OnStop(func(ctx context.Context, s Service) error { return s.(*service).Stop(ctx) }),
		),
	),

	fx.Decorate(
		func(logger *zap.Logger) *zap.Logger { return logger.Named(ModuleName) },
	),
)
