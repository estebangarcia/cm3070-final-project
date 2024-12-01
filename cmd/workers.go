package cmd

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/workers"
	"github.com/spf13/cobra"
)

// workersCmd represents the server command
var workersCmd = &cobra.Command{
	Use:   "workers",
	Short: "Start registry background workers",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg config.AppConfig

		if err := env.Parse(&cfg); err != nil {
			log.Fatal(err)
		}

		ctx := NewSigKillContext()

		sqsClient := helpers.GetSQSClient(ctx, cfg)
		dbClient, err := helpers.GetDBClient(ctx, &cfg)
		if err != nil {
			log.Fatal(err)
		}

		userRepository := repositories.NewUserRepository(dbClient)

		userSignedUpWorker := workers.NewUserSignupWorker(sqsClient, userRepository)

		sqsWorkerDispatcher := workers.NewSQSWorkerDispatcher(cfg.SignupWorker.QueueURL, sqsClient, 10)
		sqsWorkerDispatcher.Start(ctx, userSignedUpWorker)
	},
}

func init() {
	rootCmd.AddCommand(workersCmd)
}
