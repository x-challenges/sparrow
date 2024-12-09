package instruments

import (
	"math/big"
)

// Instrument
type Instrument struct {
	Address  string `json:"address" validate:"required"`
	Ticker   string `json:"ticker" validate:"required"`
	Decimals int    `json:"decimals" validate:"required"`
	zeros    int64
}

// QFromInt64
func (i *Instrument) QFromInt64(amount int64) int64 { return amount * i.zeros }

// QFromBigFloat
func (i *Instrument) QFromBigFloat(amount *big.Float) int64 {
	var res, _ = amount.Mul(amount, new(big.Float).SetInt64(i.zeros)).Int64()
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
