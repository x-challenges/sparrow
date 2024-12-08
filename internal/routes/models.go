package routes

import (
	"iter"
	"math/big"

	"sparrow/internal/instruments"
)

// Iterator
type Iterator iter.Seq[*Route]

// Route
type Route struct {
	nextRoute *Route
	// prevRoute *Route

	// instruments
	Base  *instruments.Instrument
	Quote *instruments.Instrument

	// other
	Amount   *big.Float
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

// Range
func (p *Pool) Range() Iterator {
	return func(yield func(*Route) bool) {
		for _, rp := range p.priority {
			var current = rp.head

			for {
				if !yield(current) {
					return
				}

				// exist if end of linked list
				if current == rp.tail {
					break
				}

				// exit if next route empty
				if current.nextRoute == nil {
					break
				}

				current = current.nextRoute
			}
		}
	}
}
