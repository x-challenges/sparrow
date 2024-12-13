package server

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alitto/pond/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"sparrow/internal/block"
	"sparrow/internal/prices"
	"sparrow/internal/quotes"
	"sparrow/internal/routes"
)

// Server
type Server struct {
	logger *zap.Logger
	quotes quotes.Service
	routes routes.Service
	blocks block.Listener
	prices prices.Service

	config *Config
	pool   pond.ResultPool[*quotes.Quotes]

	stop chan struct{}
	done *sync.WaitGroup
}

// NewServerFx
type NewServerFx struct {
	fx.In

	Logger *zap.Logger
	Quotes quotes.Service
	Routes routes.Service
	Blocks block.Service
	Prices prices.Service

	Config *Config
}

// NewServer
func NewServer(p NewServerFx) *Server {
	return &Server{
		logger: p.Logger,
		quotes: p.Quotes,
		routes: p.Routes,
		blocks: p.Blocks.Subscribe(),
		prices: p.Prices,

		config: p.Config,
		pool:   pond.NewResultPool[*quotes.Quotes](p.Config.Server.Concurrency),

		stop: make(chan struct{}),
		done: &sync.WaitGroup{},
	}
}

// Start
func (s *Server) Start(ctx context.Context) error {
	s.done.Add(1)

	go func() {
		defer s.done.Done()
		defer s.pool.StopAndWait()

		var processCancel context.CancelFunc

		for {
			select {
			case <-s.stop:
				return

			case block := <-s.blocks.Updates():
				// cancel previous tasks
				if processCancel != nil {
					processCancel()
				}

				// create new context with cancelation
				ctx, processCancel = context.WithCancel(context.Background())

				// run process goroutine in background
				go func() {
					defer processCancel()
					s.Process(ctx, block)
				}()
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
		started          = time.Now()
		counter          int32
		groupCtx, cancel = context.WithTimeout(ctx, s.config.Server.Deadline)
		group            = s.pool.NewGroupContext(groupCtx)
	)
	defer cancel()

	// prepare logger
	var logger = s.logger.With(zap.String("block.id", block.ID))

	// iterate by available routes
	for route := range s.routes.Range() {
		var route = *route

		// launch task in a group
		group.Submit(
			func() *quotes.Quotes {
				atomic.AddInt32(&counter, 1)

				var (
					input  = route.Base
					output = route.Quote
					amount = route.Base.QFromBigFloat(route.Amount)
				)

				var (
					qs  *quotes.Quotes
					err error
				)

				// enrich logs
				var logger = logger.With(
					zap.String("input", input.Address),
					zap.String("output", output.Address),
					zap.Int64("amount", amount),
				)

				// take quotes
				if qs, err = s.quotes.Quotes(groupCtx, input, output, amount); err != nil {
					if !errors.Is(err, quotes.ErrQuoteNotFound) {
						logger.Error("take quotes failed", zap.Error(err))
					}
					return nil
				}

				// check profit
				if !qs.HasProfit() {
					return nil
				}

				logger.Info("quotes has profit",
					zap.Any("quotes", qs),
				)

				return qs
			},
		)
	}

	// wait group
	var (
		stats []*quotes.Quotes
		err   error
	)

	if stats, err = group.Wait(); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			logger.Info("group task canceled")

		case errors.Is(err, context.DeadlineExceeded):
			logger.Info("group task canceled by timeout")

		default:
			logger.Error("group task failed", zap.Error(err))
		}
	}

	// save stats
	go func(stats ...*quotes.Quotes) {
		if err = s.quotes.BatchInsert(context.Background(), stats...); err != nil {
			logger.Error("save stats failed", zap.Error(err))
		}
	}(stats...)

	logger.Info("process completed",
		zap.Int32("routes", counter),
		zap.Duration("elapsed", time.Since(started)),
	)
}
