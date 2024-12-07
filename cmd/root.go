package cmd

import (
	"github.com/spf13/cobra"

	"github.com/x-challenges/raven/modules/config"
)

var rootCmd = &cobra.Command{}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().
		StringSliceVar(&config.Files, "config", nil, "configuration files: yaml, toml and json format")
}
