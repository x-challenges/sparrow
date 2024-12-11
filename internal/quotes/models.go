package quotes

import (
	"database/sql/driver"

	"github.com/x-challenges/raven/kun/model"

	"sparrow/internal/jupiter"
)

// Quote
type Quote struct {
	Quote     *jupiter.Quote `json:"quote"`
	StartedAt int64          `json:"started_at"`
	EndedAt   int64          `json:"ended_at"`
}

// Quotes
type Quotes struct {
	Direct  *Quote `json:"direct"`
	Reverse *Quote `json:"reverse"`
}

func (q *Quotes) Scan(src interface{}) error  { return model.JSONScanner(q, src) }
func (q Quotes) Value() (driver.Value, error) { return model.JSONValuer(q) }

// Profit
func (q *Quotes) Profit() (float32, bool) {
	var yes = q.Direct.Quote.OutAmount > q.Reverse.Quote.InAmount

	if yes {
		return (1.0 - float32(q.Reverse.Quote.InAmount)/float32(q.Direct.Quote.OutAmount)) * 100.0, true
	}

	return 0, false
}

// Model
type Model struct {
	model.Base `gorm:"embedded"`
	Quotes     Quotes `gorm:"column:quotes"`
}

// TableName
func (Model) TableName() string { return "quotes_stat" }
