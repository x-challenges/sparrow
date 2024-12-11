package tokens

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ModuleName
const ModuleName = "tokens"

// Module
var Module = fx.Module(
	ModuleName,

	// public usage
	fx.Provide(
		fx.Annotate(
			newService,
			fx.As(new(Service)),
			fx.OnStart(func(ctx context.Context, s Service) error { return s.(*service).Start(ctx) }),
		),
	),

	// private usage
	fx.Provide(
		fx.Private,
		fx.Annotate(newRepository, fx.As(new(Repository))),
	),

	fx.Decorate(
		func(logger *zap.Logger) *zap.Logger { return logger.Named(ModuleName) },
	),
)
