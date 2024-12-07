package instruments

// Status
type Status = string

const (
	Active   Status = "ACTIVE"
	InActive Status = "INACTIVE"
)

// Instrument
type Instrument struct {
	Ticker  string `json:"ticker" validate:"required"`
	Status  Status `json:"status" validate:"required"`
	Address string `json:"address" validate:"required"`
}

// NewInstrument
func NewInstrument(ticker, address string) *Instrument {
	return &Instrument{
		Ticker:  ticker,
		Status:  Active,
		Address: address,
	}
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
