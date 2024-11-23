package config

import "fmt"

type CognitoConfig struct {
	Url          string `env:"POOL_URL,notEmpty"`
	ClientId     string `env:"CLIENT_ID,notEmpty"`
	ClientSecret string `env:"CLIENT_SECRET,notEmpty"`
	Region       string `env:"REGION,notEmpty" envDefault:"eu-west-1"`
}

type AppConfig struct {
	ServerPort uint16        `env:"SERVER_PORT,notEmpty" envDefault:"8081"`
	BaseURL    string        `env:"BASE_URL,notEmpty"`
	Cognito    CognitoConfig `envPrefix:"COGNITO_"`
}

func (a AppConfig) GetCognitoJWKUrl() string {
	return fmt.Sprintf("%s/.well-known/jwks.json", a.Cognito.Url)
}
