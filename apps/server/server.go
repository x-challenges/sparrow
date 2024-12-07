package server

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"sparrow/internal/jupyter"
	"sparrow/internal/routes"
)

// Server
type Server struct {
	logger  *zap.Logger
	jupyter jupyter.Client
	routes  routes.Service

	config *Config

	stop chan struct{}
	done *sync.WaitGroup
}

// NewServer
func NewServer(logger *zap.Logger, jupyter jupyter.Client, routes routes.Service, config *Config) *Server {
	return &Server{
		logger:  logger,
		jupyter: jupyter,
		routes:  routes,

		config: config,

		stop: make(chan struct{}),
		done: &sync.WaitGroup{},
	}
}

// Start
func (s *Server) Start(ctx context.Context) error {
	s.done.Add(1)

	go func() {
		defer s.done.Done()

		var ticker = time.NewTicker(s.config.Server.Ticker)

		for {
			select {
			case <-s.stop:
				return

			case t := <-ticker.C:
				s.logger.Info("tick", zap.Time("time", t))

				s.Process(ctx)
			}
		}

	}()

	return nil
}

// Stop
func (s *Server) Stop(context.Context) error {
	// send stop signal
	close(s.stop)

	// wait done
	s.done.Wait()

	// exit
	return nil
}

// Process
func (s *Server) Process(context.Context) {
	for route := range s.routes.Range(context.Background()) {
		s.logger.Info("route",
			zap.String("base", route.Base.Ticker),
			zap.String("quote", route.Quote.Ticker),
			zap.String("amount", route.Amount.FloatString(5)),
		)
	}
}
