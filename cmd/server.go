package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/x-challenges/raven/modules/config"
	"github.com/x-challenges/raven/modules/fasthttp"
	"github.com/x-challenges/raven/modules/gorm"
	"github.com/x-challenges/raven/modules/logger"
	"github.com/x-challenges/raven/modules/yandex/ydb"

	"sparrow/apps/server"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(_ *cobra.Command, _ []string) {
		app := fx.New(
			fx.RecoverFromPanics(),

			// raven
			config.Module(),
			logger.Module,
			fasthttp.Module,
			gorm.Module,
			ydb.Module,

			// apps
			server.Module,
		)

		// run app
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
