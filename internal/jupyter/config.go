package jupyter

// Config
type Config struct {
	Jupyter struct {
		Price struct {
			APIHost string `mapstructure:"api_host" validate:"required"`
		} `mapstructure:"price"`

		Quote struct {
			APIHost string `mapstructure:"api_host" validate:"required"`
		} `mapstructure:"quote"`

		Token struct {
			APIHost string `mapstructure:"api_host" validate:"required"`
		} `mapstructure:"token"`
	} `mapstructure:"jupyter"`
}
