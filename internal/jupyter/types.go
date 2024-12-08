package jupyter

// Token
type Token struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int64  `json:"decimals"`
	LogoURI  string `json:"logoURI"`
}

// QuoteRoutePlanSwapInfo
type QuoteRoutePlanSwapInfo struct {
	Ammkey     string `json:"ammKey"`
	Label      string `json:"label"`
	InputMint  string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	InAmount   int64  `json:"inAmount,string"`
	OutAmount  int64  `json:"outAmount,string"`
	FeeAmount  int64  `json:"feeAmount"`
}

// QuoteRoutePlan
type QuoteRoutePlan struct {
	SwapInfo any `json:"swapInfo"`
	Percent  int `json:"percent"`
}

// Quote
type Quote struct {
	InputMint  string           `json:"inputMint"`
	InAmount   int64            `json:"inAmont,string"`
	OutputMint string           `json:"outputMint"`
	OutAmount  int64            `json:"outAmount,string"`
	RoutePlan  []QuoteRoutePlan `json:"routePlan"`
	TimeTaken  float32          `json:"timeTaken"`
}

// Quotes
type Quotes struct {
	Direct  *Quote `json:"direct"`
	Reverse *Quote `json:"reverse"`
}

// HasProfit
func (q *Quotes) HasProfit() bool {
	return q.Direct.InAmount < q.Reverse.OutAmount
}
