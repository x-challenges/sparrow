package jupiter

// QuoteOptions
type QuoteOptions struct {
	SwapMode SwapMode
}

// NewQuoteOptions
func NewQuoteOptions() *QuoteOptions {
	return &QuoteOptions{
		SwapMode: ExactIn,
	}
}

// Apply
func (qo *QuoteOptions) Apply(options ...QuoteOptionFunc) *QuoteOptions {
	for _, option := range options {
		option(qo)
	}

	return qo
}

// QuoteOptionFunc
type QuoteOptionFunc func(*QuoteOptions)

// WithSwapMode
func WithSwapMode(sm SwapMode) QuoteOptionFunc {
	return func(qo *QuoteOptions) {
		qo.SwapMode = sm
	}
}
