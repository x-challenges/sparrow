package jupyter

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ModuleName
const ModuleName = "jupyter"

// Module
var Module = fx.Module(
	ModuleName,

	// public usage
	fx.Provide(
		fx.Annotate(newClient, fx.As(new(Client))),
	),

	fx.Decorate(
		func(logger *zap.Logger) *zap.Logger { return logger.Named(ModuleName) },
	),
)
