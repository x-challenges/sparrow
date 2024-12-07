package server

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

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

	// force
	fx.Invoke(func(routes.Service) {}),

	fx.Decorate(
		func(logger *zap.Logger) *zap.Logger { return logger.Named(ModuleName) },
	),
)
