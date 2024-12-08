package jupyter

import (
	"sync/atomic"
)

// Balancer
type Balancer struct {
	servers []string
	index   int32
}

// NewBalancer
func NewBalancer(servers ...string) *Balancer {
	return &Balancer{
		servers: servers,
		index:   0,
	}
}

// Next
func (b *Balancer) Next() string {
	var index = atomic.SwapInt32(&b.index, b.index%int32(len(b.servers))) // nolint:gosec

	return b.servers[index]
}
