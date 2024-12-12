package instruments

import (
	"iter"
	"math/big"
)

// Tag
type Tag = string

const (
	Unspecified Tag = ""
	Base        Tag = "base"
	Route       Tag = "route"
	Swap        Tag = "swap"
)

// Iterator
type Iterator iter.Seq[*Instrument]

// Instrument
type Instrument struct {
	Address    string     `json:"address"`
	Ticker     string     `json:"ticker"`
	Decimals   int        `json:"decimals"`
	Tags       []Tag      `json:"tags"`
	Zeros      int64      `json:"-"`
	token      *Token     `json:"-"`
	zerosValue *big.Float `json:"-"`
}

// QFromInt64
func (i *Instrument) QFromInt64(amount int64) int64 { return amount * i.Zeros }

// FFromInt64
func (i *Instrument) FFromInt64(amount int64) float64 { return float64(amount) / float64(i.Zeros) }

// QFromBigFloat
func (i *Instrument) QFromBigFloat(amount *big.Float) int64 {
	var res, _ = new(big.Float).Mul(amount, i.zerosValue).Int64()
	return res
}

// Instruments
type Instruments []*Instrument

// Addresses
func (is Instruments) Addresses() []string {
	var result = make([]string, len(is))

	for _, i := range is {
		result = append(result, i.Address)
	}

	return result
}
