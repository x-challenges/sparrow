package routes

// Config
type Config struct {
	Routes struct {
		BaseCcy []string `mapstructure:"base_ccy" validate:"required"`
		Range   [2]int64 `mapstructure:"range" validate:"required" default:"[1,10]"`
		Step    int64    `mapstructure:"step" validate:"required" default:"1000"`

		// Overrides
		Overrides []struct {
			Instrument string    `mapstructure:"instrument" validate:"required"`
			Priority   *int      `mapstructure:"priority"`
			Range      *[2]int64 `mapstructure:"range"`
			Step       *int64    `mapstructure:"step"`
		} `mapstructure:"pool"`
	} `mapstructure:"routes"`
}
