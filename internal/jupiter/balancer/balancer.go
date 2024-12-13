package balancer

// Balancer
type Balancer interface {
	// String
	Next() string
}

// NewBalancer
func NewBalancer(servers ...string) Balancer {
	if len(servers) == 1 {
		return NewFirst(servers...)
	}
	return NewRoundRobin(servers...)
}
