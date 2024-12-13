package balancer

import (
	"sync"
)

// Balancer
type RoundRobin struct {
	servers []string
	index   int32
	mu      *sync.Mutex
}

var _ Balancer = (*RoundRobin)(nil)

// NewRoundRobin
func NewRoundRobin(servers ...string) *RoundRobin {
	return &RoundRobin{
		servers: servers,
		index:   0,
		mu:      &sync.Mutex{},
	}
}

// Next
func (b *RoundRobin) Next() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.index = (b.index + 1) % int32(len(b.servers)) // nolint:gosec

	return b.servers[b.index]
}
