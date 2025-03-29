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

// This middleware extracts base64 encoded authentication credentials
// it then verifies them against cognito and stores the token in the context
// to be used further down the middleware pipeline
func (a *ExtractBasicCredentialsMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		// If the header is missing then return 401
		if header == "" {
			w.Header().Set(wwwAuthenticateHeader, "Basic")
			responses.OCIUnauthorizedError(w)
			return
		}

		// If the header specifies the Bearer schema then continue the request
		// and let the token middleware handle this
		if strings.HasPrefix(header, bearerSchema) {
			next.ServeHTTP(w, r.WithContext(r.Context()))
			return
		}

		// If the header doesn't contain the Basic schema then return 401
		if !strings.HasPrefix(header, basicSchema) {
			w.Header().Set(wwwAuthenticateHeader, "Basic")
			responses.OCIUnauthorizedError(w)
			return
		}

		// Extract and decode the credentials from the header
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
			w.Header().Set(wwwAuthenticateHeader, "Basic")
			responses.OCIUnauthorizedError(w)
			return
		}

		// Verify that the credentials are valid

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

		// Store JWT token in context to be validated by the token middleware
		ctx := context.WithValue(r.Context(), "token", *output.AuthenticationResult.AccessToken)
		ctx = context.WithValue(ctx, "expires_in", output.AuthenticationResult.ExpiresIn)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
