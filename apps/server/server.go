package server

import (
	"context"
	"errors"
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
	cancel context.CancelFunc

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
		pool:   pond.NewPool(p.Config.Server.Concurrency),

		stop: make(chan struct{}),
		done: &sync.WaitGroup{},
	}
}

// Start
func (s *Server) Start(context.Context) error {
	s.done.Add(1)

	go func() {
		defer s.done.Done()
		defer s.pool.StopAndWait()

		// var processCancel context.CancelFunc

		for {
			select {
			case <-s.stop:
				return

			case block := <-s.blocks.Updates():
				// cancel previous tasks
				// if processCancel != nil {
				// 	processCancel()
				// }

				// create new context with cancelation
				// ctx, processCancel = context.WithCancel(ctx)

				// run porcess goroutine in background
				// go func() {
				// 	// defer processCancel()
				// 	s.Process(ctx, block)
				// }()

				s.Process(context.Background(), block)
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
func (s *Server) Process(ctx context.Context, block *block.Block) {
	var (
		groupCtx, cancel = context.WithTimeout(ctx, s.config.Server.Deadline)
		group            = s.pool.NewGroupContext(groupCtx)
		err              error
	)
	defer cancel()

	// iterate by available routes
	for route := range s.routes.Range() {
		group.Submit(
			func() {
				var quotes *jupyter.Quotes

				var (
					input  = route.Base
					output = route.Quote
					amount = route.Base.QFromBigFloat(route.Amount)
				)

				// take quotes
				if quotes, err = s.jupyter.Quotes(groupCtx, input, output, amount); err != nil {
					s.logger.Error("take quotes failed", zap.Error(err))
				}

				// check profit
				if profit, yes := quotes.Profit(); yes {
					s.logger.Info("quotes has profit", zap.String("block.id", block.ID), zap.Float32("profit", profit))
				}
			},
		)
	}

	// wait group
	if err = group.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			s.logger.Info("group task canceled", zap.String("block.id", block.ID))
			return
		}

		s.logger.Error("group task failed", zap.Error(err))
	}
}
