package jupiter

import (
	"sync"
)

// Balancer
type Balancer struct {
	servers []string
	index   int32
	mu      *sync.Mutex
}

// NewBalancer
func NewBalancer(servers ...string) *Balancer {
	return &Balancer{
		servers: servers,
		index:   0,
		mu:      &sync.Mutex{},
	}
}

// Next
func (b *Balancer) Next() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.index = (b.index + 1) % int32(len(b.servers)) // nolint:gosec

	return b.servers[b.index]
}
