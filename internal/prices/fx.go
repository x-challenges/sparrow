package prices

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ModuleName
const ModuleName = "prices"

// Module
var Module = fx.Module(
	ModuleName,

	// public usage
	fx.Provide(
		fx.Annotate(newService, fx.As(new(Service))),
	),

	// private usage
	fx.Provide(
		fx.Private,
		NewRates,
	),

	fx.Decorate(
		func(logger *zap.Logger) *zap.Logger { return logger.Named(ModuleName) },
	),
)
