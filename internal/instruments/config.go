package instruments

// Config
type Config struct {
	Instruments struct {
		Pool []Instrument `mapstructure:"pool" validate:"required"`
	} `mapstructure:"instruments"`
}
