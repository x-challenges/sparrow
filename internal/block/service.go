package block

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

// Service
type Service interface {
	// Subscribe
	Subscribe() Listener
}

// Service interface implementation
type service struct {
	logger       *zap.Logger
	producer     Producer
	subscription Subscription

	stop chan struct{}
	done *sync.WaitGroup
}

var _ Service = (*service)(nil)

// NewService
func newService(logger *zap.Logger, producer Producer) (*service, error) {
	return &service{
		logger:       logger,
		producer:     producer,
		subscription: NewSubscription(),

		stop: make(chan struct{}),
		done: &sync.WaitGroup{},
	}, nil
}

// Start
func (s *service) Start(ctx context.Context) error {
	// launch background subscription
	s.done.Add(1)
	go func() {
		defer s.done.Done()

		s.subscription.Start(ctx)
	}()

	// launch background producer
	s.done.Add(1)
	go func() {
		defer s.done.Done()

		if err := s.producer.Start(ctx); err != nil {
			s.logger.Fatal("producer start failed", zap.Error(err))
		}
	}()

	// broadcast blocks to subscribers
	s.done.Add(1)
	go func() {
		defer s.done.Done()

		for {
			select {
			case <-s.stop:
				return
			case block := <-s.producer.Channel():
				s.subscription.Publish(block)
			}
		}
	}()

	return nil
}

// Stop
func (s *service) Stop(ctx context.Context) error {
	// send close signal
	close(s.stop)

	// close producer
	_ = s.producer.Stop(ctx)

	// close subscription
	s.subscription.Stop(ctx)

	// wait done
	s.done.Wait()

	// exit
	return nil
}

// Subscribe implements Service interface
func (s *service) Subscribe() Listener {
	return s.subscription.Subscribe()
}
