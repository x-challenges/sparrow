package orders

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ModuleName
const ModuleName = "orders"

var Module = fx.Module(
	ModuleName,

	// public usage
	fx.Provide(
		fx.Annotate(newService, fx.As(new(Service))),
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
