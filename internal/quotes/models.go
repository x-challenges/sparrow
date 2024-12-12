package quotes

import (
	"database/sql/driver"
	"time"

	"github.com/x-challenges/raven/kun/model"

	"sparrow/internal/instruments"
	"sparrow/internal/jupiter"
)

// Quote
type Quote = jupiter.Quote

// Quotes
type Quotes struct {
	Direct  *Quote                  `json:"direct"`
	Reverse *Quote                  `json:"reverse"`
	Input   *instruments.Instrument `json:"-"`
	Output  *instruments.Instrument `json:"-"`
	Diff    int64                   `json:"diff"`
	Profit  float32                 `json:"profit"`
	Elapsed time.Duration           `json:"elapsed"`
}

func (q *Quotes) Scan(src interface{}) error  { return model.JSONScanner(q, src) }
func (q Quotes) Value() (driver.Value, error) { return model.JSONValuer(q) }

// HasProfit
func (q *Quotes) HasProfit() bool { return q.Diff > 0 }

// Model
type Model struct {
	model.Base `gorm:"embedded"`
	Quotes     Quotes `gorm:"column:quotes"`
}

// TableName
func (Model) TableName() string { return "quotes_stat" }
