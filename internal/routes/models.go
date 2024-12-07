package routes

import (
	"math/big"

	"sparrow/internal/instruments"
)

// Route
type Route struct {
	prevRoute *Route
	nextRoute *Route

	// instruments
	Base  *instruments.Instrument
	Quote *instruments.Instrument

	// other
	Amount   *big.Rat
	Priority int
}

// RoutePoint
type RoutePoint struct {
	head *Route
	tail *Route
}

// Insert
func (pp *RoutePoint) Insert(route *Route) {
	if pp.head == nil {
		pp.head = route
		pp.tail = route
	} else {
		pp.tail.nextRoute = route
		pp.tail = route
	}
}

// Pool
type Pool struct {
	index    []*Route
	priority map[int]*RoutePoint
}

// NewPool
func NewPool() *Pool {
	return &Pool{
		index:    make([]*Route, 0, 1000),
		priority: make(map[int]*RoutePoint, 1000),
	}
}

// AddRoute
func (p *Pool) AddRoute(route *Route) {
	var (
		rp    *RoutePoint
		exist bool
	)

	// allocate new route point if not exists
	if rp, exist = p.priority[route.Priority]; !exist {
		rp = new(RoutePoint)
		p.priority[route.Priority] = rp
	}

	// insert route to route point
	rp.Insert(route)

	// add route to index
	p.index = append(p.index, route)
}
