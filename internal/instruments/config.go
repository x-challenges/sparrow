package instruments

// Config
type Config struct {
	Instruments struct {
		// Loader settings for tokens
		Loader struct {
			Skip struct {
				DailyVolume float64 `mapstructure:"daily_volume" default:"1_000_000"`
			} `mapstructure:"skip"`
		} `mapstructure:"loader"`

		// Available instruments
		Pool []struct {
			Address string `mapstructure:"address" validate:"required"`
			Tags    []Tag  `mapstructure:"tags" validate:"required" default:"[swap]"`
		} `mapstructure:"pool" validate:"required"`
	} `mapstructure:"instruments"`
}
