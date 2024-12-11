package prices

// Config
type Config struct {
	Prices struct {
		Loader struct {
			ChunkSize int `mapstructure:"chunk_size" validate:"required,gt=0,lte=100" default:"100"`
		} `mapstructure:"loader"`
	} `mapstructure:"prices"`
}
