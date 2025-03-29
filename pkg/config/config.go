package config

import "fmt"

type AdminUserConfig struct {
	// Admin user email
	Email string `env:"EMAIL,notEmpty"`
	// Admin user password
	Password string `env:"PASSWORD,notEmpty"`
	// Admin user SUB ID
	Sub string `env:"SUB,notEmpty"`
}

type DatabaseConfig struct {
	// DSN to connect to
	DSN string `env:"DSN,notEmpty"`
	// Is the debug feature enabled
	Debug bool `env:"DEBUG_MODE" envDefault:"false"`
}

type CognitoConfig struct {
	// URL for the Cognito Pool
	Url string `env:"POOL_URL,notEmpty"`
	// ID for the client
	ClientId string `env:"CLIENT_ID,notEmpty"`
	// Client secret to be used for validating credentials
	ClientSecret string `env:"CLIENT_SECRET,notEmpty"`
	// Region where Cognito is deployed
	Region string `env:"REGION,notEmpty" envDefault:"eu-west-1"`
}

type S3Config struct {
	// S3 bucket name where we store uploads
	BlobsBucketName string `env:"BLOBS_BUCKET_NAME,notEmpty" envDefault:"egarcia-blob-uploads"`
}

type SESConfig struct {
	// Email used to send emails
	FromEmailAddress string `env:"FROM_EMAIL,notEmpty" envDefault:"elg4@student.london.ac.uk"`
}

type SignupWorkerConfig struct {
	// SQS Queue URL where Cognito will send new signups
	QueueURL string `env:"QUEUE_URL,notEmpty" envDefault:"https://sqs.eu-west-1.amazonaws.com/205930648580/user-signed"`
}

type AppConfig struct {
	// Port the server will be run on
	ServerPort uint16 `env:"SERVER_PORT,notEmpty" envDefault:"8081"`
	// Base URL for the API to be used for redirections
	BaseURL string `env:"BASE_URL,notEmpty"`
	// Base URL for the frontend UI to be used for emails
	FrontendBaseURL string `env:"FRONTEND_BASE_URL" envDefault:"http://localhost:3000"`
	// Configuration for admin user to be used for internal requests
	AdminUser AdminUserConfig `envPrefix:"ADMIN_"`
	// Database configuration
	Database DatabaseConfig `envPrefix:"DB_"`
	// Cognito Auth configuration
	Cognito CognitoConfig `envPrefix:"COGNITO_"`
	// S3 configuration
	S3 S3Config `envPrefix:"S3_"`
	// SES Email service configuration
	SES SESConfig `envPrefix:"SES_"`
	// Configuration for the signup worker
	SignupWorker SignupWorkerConfig `envPrefix:"SIGNUP_WORKER_"`
	// Minimum size for a chunk when receiving blobs
	ChunkMinLength uint32 `env:"CHUNK_MIN_LENGTH" envDefault:"5242880"`
	// Length of the buffer for chunking in memory
	ChunkBufferLength uint32 `env:"CHUNK_BUFFER_LENGTH" envDefault:"52428800"`
	// Go routine limit for chunking blobs in memory
	BlobUploadMaxGoRoutines int `env:"BLOB_UPLOAD_MAX_GO_ROUTINES" envDefault:"10"`
}

func (a AppConfig) GetCognitoJWKUrl() string {
	return fmt.Sprintf("%s/.well-known/jwks.json", a.Cognito.Url)
}

func (a AppConfig) GetBaseUrl() string {
	return fmt.Sprintf("https://%s", a.BaseURL)
}
