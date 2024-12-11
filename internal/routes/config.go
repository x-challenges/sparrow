package routes

// Config
type Config struct {
	Routes struct {
		Range [2]int64 `mapstructure:"range" validate:"required" default:"[1,10]"`
		Step  int64    `mapstructure:"step" validate:"required" default:"1000"`
	} `mapstructure:"routes"`
}
