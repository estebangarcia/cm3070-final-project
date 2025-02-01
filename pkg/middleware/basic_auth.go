package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

const basicSchema = "Basic "

type ExtractBasicCredentialsMiddleware struct {
	Config        *config.AppConfig
	CognitoClient *cognitoidentityprovider.Client
}

func (a *ExtractBasicCredentialsMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			responses.OCIUnauthorizedError(w)
			return
		}

		if strings.HasPrefix(header, bearerSchema) {
			next.ServeHTTP(w, r.WithContext(r.Context()))
			return
		}

		if !strings.HasPrefix(header, basicSchema) {
			responses.OCIUnauthorizedError(w)
			return
		}

		basicCreds := header[len(basicSchema):]

		decoded, err := base64.StdEncoding.DecodeString(basicCreds)
		if err != nil {
			responses.OCIUnauthorizedError(w)
			return
		}

		usernamePassword := strings.Split(string(decoded), ":")
		if len(usernamePassword) < 2 {
			responses.OCIUnauthorizedError(w)
			return
		}

		username := usernamePassword[0]
		password := usernamePassword[1]

		if username == "" || password == "" {
			responses.OCIUnauthorizedError(w)
			return
		}

		mac := hmac.New(sha256.New, []byte(a.Config.Cognito.ClientSecret))
		mac.Write([]byte(username + a.Config.Cognito.ClientId))

		secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		params := map[string]string{
			"USERNAME":    username,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		}

		// Initiate authentication request
		input := &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow:       types.AuthFlowTypeUserPasswordAuth,
			ClientId:       aws.String(a.Config.Cognito.ClientId),
			AuthParameters: params,
		}

		// Call Cognito
		output, err := a.CognitoClient.InitiateAuth(r.Context(), input)
		if err != nil {
			responses.OCIUnauthorizedError(w)
			return
		}

		ctx := context.WithValue(r.Context(), "token", *output.AuthenticationResult.AccessToken)
		ctx = context.WithValue(ctx, "expires_in", output.AuthenticationResult.ExpiresIn)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
