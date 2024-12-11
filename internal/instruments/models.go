package instruments

import (
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

// Instrument
type Instrument struct {
	Address    string `mapstructure:"address" json:"address" validate:"required"`
	Ticker     string `mapstructure:"ticker" json:"ticker" validate:"required"`
	Decimals   int    `mapstructure:"decimals" json:"decimals" validate:"required"`
	Tags       []Tag  `mapstructure:"tags" json:"tags" validate:"required" default:"[swap]"`
	Zeros      int64
	zerosValue *big.Float
}

// QFromInt64
func (i *Instrument) QFromInt64(amount int64) int64 { return amount * i.Zeros }

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
