package config

import "fmt"

type DatabaseConfig struct {
	DSN string `env:"DSN,notEmpty"`
}

type CognitoConfig struct {
	Url          string `env:"POOL_URL,notEmpty"`
	ClientId     string `env:"CLIENT_ID,notEmpty"`
	ClientSecret string `env:"CLIENT_SECRET,notEmpty"`
	Region       string `env:"REGION,notEmpty" envDefault:"eu-west-1"`
}

type S3Config struct {
	BlobsBucketName string `env:"BLOBS_BUCKET_NAME,notEmpty" envDefault:"egarcia-blob-uploads"`
}

type SignupWorkerConfig struct {
	QueueURL string `env:"QUEUE_URL,notEmpty" envDefault:"https://sqs.eu-west-1.amazonaws.com/205930648580/user-signed"`
}

type AppConfig struct {
	ServerPort              uint16             `env:"SERVER_PORT,notEmpty" envDefault:"8081"`
	BaseURL                 string             `env:"BASE_URL,notEmpty"`
	Database                DatabaseConfig     `envPrefix:"DB_"`
	Cognito                 CognitoConfig      `envPrefix:"COGNITO_"`
	S3                      S3Config           `envPrefix:"S3_"`
	SignupWorker            SignupWorkerConfig `envPrefix:"SIGNUP_WORKER_"`
	ChunkMinLength          uint32             `env:"CHUNK_MIN_LENGTH" envDefault:"5242880"`
	ChunkBufferLength       uint32             `env:"CHUNK_BUFFER_LENGTH" envDefault:"52428800"`
	BlobUploadMaxGoRoutines int                `env:"BLOB_UPLOAD_MAX_GO_ROUTINES" envDefault:"10"`
}

func (a AppConfig) GetCognitoJWKUrl() string {
	return fmt.Sprintf("%s/.well-known/jwks.json", a.Cognito.Url)
}
