package jupiter

// SwapMode
type SwapMode int

const (
	ExactIn SwapMode = iota
	ExactOut
)

// Token
type Token struct {
	Address     string  `json:"address"`
	Name        string  `json:"name"`
	Symbol      string  `json:"symbol"`
	Decimals    int     `json:"decimals"`
	DailyVolume float64 `json:"daily_volume"`
}

// Tokens
type Tokens []Token

// QuoteRoutePlanSwapInfo
type QuoteRoutePlanSwapInfo struct {
	// Ammkey     string `json:"ammKey"`
	// Label      string `json:"label"`
	// InputMint  string `json:"inputMint"`
	// OutputMint string `json:"outputMint"`
	// InAmount  int64  `json:"inAmount,string"`
	// OutAmount int64  `json:"outAmount,string"`
	FeeAmount int64  `json:"feeAmount,string"`
	FeeMint   string `json:"feeMint"`
}

// QuoteRoutePlan
type QuoteRoutePlan struct {
	SwapInfo QuoteRoutePlanSwapInfo `json:"swapInfo"`
	Percent  int                    `json:"percent"`
}

// Quote
type Quote struct {
	InputMint  string           `json:"inputMint"`
	InAmount   int64            `json:"inAmount,string"`
	OutputMint string           `json:"outputMint"`
	OutAmount  int64            `json:"outAmount,string"`
	RoutePlan  []QuoteRoutePlan `json:"routePlan"`
	TimeTaken  float32          `json:"timeTaken"`
}

// Price
type Price struct {
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	Price float64 `json:"price,string"`
}

// Prices
type Prices struct {
	Data      map[string]Price `json:"data"`
	TimeTaken float32          `json:"timeTaken"`
}
