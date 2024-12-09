package quotes

// Config
type Config struct {
	Quotes struct {
		Stats struct {
			BufferSize int `mapstructure:"buffer_size" default:"10"`
		} `mapstructure:"stats"`
	} `mapstructure:"quotes"`
}
