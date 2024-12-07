package server

import "time"

// Config
type Config struct {
	Server struct {
		Ticker time.Duration `mapstructure:"ticker" default:"1s"`
	} `mapstructure:"server"`
}
