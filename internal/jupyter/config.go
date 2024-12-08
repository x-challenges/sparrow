package jupyter

// Config
type Config struct {
	Jupyter struct {
		Token struct {
			Host string `mapstructure:"host" validate:"required"`
		} `mapstructure:"token"`

		Price struct {
			Host string `mapstructure:"host" validate:"required"`
		} `mapstructure:"price"`

		Quote struct {
			Hosts []string `mapstructure:"hosts" validate:"required"`
		} `mapstructure:"quote"`
	} `mapstructure:"jupyter"`
}
