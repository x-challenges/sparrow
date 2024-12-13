package balancer

import (
	"sync"
)

// Balancer
type RoundRobin struct {
	servers []string
	index   int32
	len     int32
	mu      *sync.Mutex
}

var _ Balancer = (*RoundRobin)(nil)

// NewRoundRobin
func NewRoundRobin(servers ...string) *RoundRobin {
	return &RoundRobin{
		servers: servers,
		index:   0,
		len:     int32(len(servers)), // nolint:gosec
		mu:      &sync.Mutex{},
	}
}

// Next
func (b *RoundRobin) Next() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.index = (b.index + 1) % b.len

	return b.servers[b.index]
}
