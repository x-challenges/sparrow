package quotes

import (
	"sparrow/internal/jupyter"
)

// Quote
type Quote = jupyter.Quote

// Quotes
type Quotes struct {
	Direct  *Quote `json:"direct"`
	Reverse *Quote `json:"reverse"`
}

// Profit
func (q *Quotes) Profit() (float32, bool) {
	var yes = q.Direct.InAmount < q.Reverse.OutAmount

	if yes {
		return (1.0 - float32(q.Direct.InAmount)/float32(q.Reverse.OutAmount)) * 100.0, true
	}

	return 0, false
}
