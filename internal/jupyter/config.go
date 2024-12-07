package jupyter

// Config
type Config struct {
	Jupyter struct {
		Price struct {
			APIHost string `mapstructure:"api_host" validate:"required" default:""`
		} `mapstructure:"price"`

		Quote struct {
			APIHost string `mapstructure:"api_host" validate:"required" default:""`
		} `mapstructure:"quote"`
	} `mapstructure:"jupyter"`
}
