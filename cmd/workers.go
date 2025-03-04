package cmd

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/workers"
	"github.com/spf13/cobra"
)

var availableWorkers = []string{
	"all",
	"security_scanner",
	"user_signup",
}

// workersCmd represents the server command
var workersCmd = &cobra.Command{
	Use:   "workers",
	Short: "Start registry background workers",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}

		if len(args) > 0 && !slices.Contains(availableWorkers, args[0]) {
			return fmt.Errorf("%s is not a valid worker", strings.ToLower(args[0]))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var cfg config.AppConfig

		if err := env.Parse(&cfg); err != nil {
			log.Fatal(err)
		}

		ctx := NewSigKillContext()

		sqsClient := helpers.GetSQSClient(ctx, cfg)
		s3Client := helpers.GetS3Client(ctx, cfg)
		dbClient, err := helpers.GetDBClient(ctx, &cfg)
		if err != nil {
			log.Fatal(err)
		}

		startWorker := "all"
		if len(args) > 0 {
			startWorker = strings.ToLower(args[0])
		}

		if startWorker == "all" || startWorker == "security_scanner" {
			fmt.Println("Starting Security Scanning worker...")
			manifestRepository := repositories.NewManifestRepository()
			periodicWorkerDispatcher := workers.NewPeriodicWorkerDispatcher(10*time.Second, dbClient)
			trivyWorker := workers.NewSecurityScannerWorker(5, s3Client, manifestRepository, &cfg)
			periodicWorkerDispatcher.Start(ctx, trivyWorker)
		}

		if startWorker == "all" || startWorker == "user_signup" {
			fmt.Println("Starting User Signup worker...")
			userRepository := repositories.NewUserRepository()
			userSignedUpWorker := workers.NewUserSignupWorker(sqsClient, userRepository)
			sqsWorkerDispatcher := workers.NewSQSWorkerDispatcher(cfg.SignupWorker.QueueURL, sqsClient, 10, dbClient)
			sqsWorkerDispatcher.Start(ctx, userSignedUpWorker)
		}

		<-ctx.Done()
	},
}

func init() {
	rootCmd.AddCommand(workersCmd)
}
