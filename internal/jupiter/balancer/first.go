package balancer

// First
type First struct {
	servers []string
}

var _ Balancer = (*First)(nil)

// NewFirst
func NewFirst(servers ...string) *First {
	return &First{
		servers: servers,
	}
}

// Next implements Balancer interface
func (d *First) Next() string { return d.servers[0] }
