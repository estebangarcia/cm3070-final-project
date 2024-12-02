package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

const bearerSchema = "Bearer "
const wwwAuthenticateHeader = "WWW-Authenticate"

type JWTAuthMiddleware struct {
	Config   *config.AppConfig
	JwkCache *jwk.Cache
}

func (a *JWTAuthMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, bearerSchema) {
			w.Header().Set(wwwAuthenticateHeader, a.getAuthenticationUrl())
			responses.OCIUnauthorizedError(w)
			return
		}

		jwkSet, err := a.JwkCache.Lookup(r.Context(), a.Config.GetCognitoJWKUrl())
		if err != nil {
			responses.OCIInternalServerError(w)
			return
		}

		jwtToken := header[len(bearerSchema):]
		token, err := jwt.Parse([]byte(jwtToken), jwt.WithKeySet(jwkSet))
		if err != nil {
			w.Header().Set(wwwAuthenticateHeader, a.getAuthenticationUrl())
			responses.OCIUnauthorizedError(w)
			return
		}

		userSub, _ := token.Subject()

		ctx := context.WithValue(r.Context(), "user_sub", userSub)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *JWTAuthMiddleware) getAuthenticationUrl() string {
	return fmt.Sprintf(`Bearer realm="%s/v2/login",service="registry.io"`, a.Config.BaseURL)
}
