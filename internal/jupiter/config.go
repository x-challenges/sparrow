package jupiter

// Config
type Config struct {
	Jupiter struct {
		Token struct {
			Host string   `mapstructure:"host" validate:"required"`
			Tags []string `mapstructure:"tags"`
		} `mapstructure:"token"`

		Price struct {
			Host string `mapstructure:"host" validate:"required"`
		} `mapstructure:"price"`

		Quote struct {
			Hosts                      []string `mapstructure:"hosts" validate:"required"`
			OnlyDirectRoutes           bool     `mapstructure:"only_direct_routes" default:"false"`
			RestrictIntermediateTokens bool     `mapstructure:"restrict_intermediate_tokens" default:"false"`
		} `mapstructure:"quote"`
	} `mapstructure:"jupiter"`
}
