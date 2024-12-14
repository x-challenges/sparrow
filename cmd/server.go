package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/x-challenges/raven/modules/config"
	"github.com/x-challenges/raven/modules/fasthttp"
	"github.com/x-challenges/raven/modules/gorm"
	"github.com/x-challenges/raven/modules/http"
	"github.com/x-challenges/raven/modules/logger"
	"github.com/x-challenges/raven/modules/monitoring"
	"github.com/x-challenges/raven/modules/worker"
	"github.com/x-challenges/raven/modules/yandex/ydb"

	"sparrow/apps/server"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(_ *cobra.Command, _ []string) {
		app := fx.New(
			// fx.NopLogger,
			fx.RecoverFromPanics(),

			// raven
			config.Module(),
			logger.Module,
			monitoring.Module,
			http.Module,
			worker.Module,
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
