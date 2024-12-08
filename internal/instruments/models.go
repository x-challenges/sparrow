package instruments

// Instrument
type Instrument struct {
	Address  string `json:"address" validate:"required"`
	Ticker   string `json:"ticker" validate:"required"`
	Decimals int64  `json:"decimals" validate:"required"`
}

// Amount
func (i *Instrument) Amount(amount int64) int64 { return amount * i.Decimals }

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
