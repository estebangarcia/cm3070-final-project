package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

const basicSchema = "Basic "

type ExtractBasicCredentialsMiddleware struct {
}

func (a *ExtractBasicCredentialsMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, basicSchema) {
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

		ctx := context.WithValue(r.Context(), "username", username)
		ctx = context.WithValue(ctx, "password", password)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
