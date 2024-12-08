package server

import "time"

// Config
type Config struct {
	Server struct {
		Concurrency int           `mapstructure:"concurrency" default:"100"`
		Ticker      time.Duration `mapstructure:"ticker" default:"1s"`
		Deadline    time.Duration `mapstructure:"deadline" default:"1s"`
	} `mapstructure:"server"`
}
