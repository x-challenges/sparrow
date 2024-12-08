package server

import (
	"context"
	"sync"

	"github.com/alitto/pond/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"sparrow/internal/block"
	"sparrow/internal/jupyter"
	"sparrow/internal/routes"
)

// Server
type Server struct {
	logger  *zap.Logger
	pool    pond.Pool
	jupyter jupyter.Client
	routes  routes.Service
	blocks  block.Listener

	config *Config

	stop chan struct{}
	done *sync.WaitGroup
}

// NewServerFx
type NewServerFx struct {
	fx.In

	Logger  *zap.Logger
	Jupyter jupyter.Client
	Routes  routes.Service
	Blocks  block.Service

	Config *Config
}

// NewServer
func NewServer(p NewServerFx) *Server {
	return &Server{
		logger:  p.Logger,
		jupyter: p.Jupyter,
		routes:  p.Routes,
		blocks:  p.Blocks.Subscribe(),

		config: p.Config,

		stop: make(chan struct{}),
		done: &sync.WaitGroup{},
	}
}

// Start
func (s *Server) Start(ctx context.Context) error {
	s.pool = pond.NewPool(s.config.Server.Concurrency,
		pond.WithContext(ctx),
	)

	s.done.Add(1)
	go func() {
		defer s.done.Done()

		for {
			select {
			case <-s.stop:
				return

			case block := <-s.blocks.Updates():
				s.Process(ctx, block)
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
func (s *Server) Process(ctx context.Context, _ *block.Block) {
	var (
		groupCtx, cancel = context.WithTimeout(ctx, s.config.Server.Deadline)
		group            = s.pool.NewGroupContext(groupCtx)
		err              error
	)
	defer cancel()

	// iterate by routes
	for range s.routes.Range(context.Background()) {
		group.Submit()
	}

	// wait group
	if err = group.Wait(); err != nil {
		s.logger.Error("group task failed", zap.Error(err))
	}
}
