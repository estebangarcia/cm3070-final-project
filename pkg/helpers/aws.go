package helpers

import (
	"context"
	"log"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
)

func GetCognitoClient(ctx context.Context, cfg config.AppConfig) *cognitoidentityprovider.Client {
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.Cognito.Region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create a Cognito Identity Provider client
	return cognitoidentityprovider.NewFromConfig(awsCfg)
}

func GetS3Client(ctx context.Context, cfg config.AppConfig) *s3.Client {
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.Cognito.Region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return s3.NewFromConfig(awsCfg)
}

func GetS3PresignClient(s3Client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(s3Client)
}
