package middleware

import (
	"fmt"
	"strings"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

const bearerSchema = "Bearer "
const wwwAuthenticateHeader = "WWW-Authenticate"

type AuthMiddleware struct {
	Config   *config.AppConfig
	JwkCache *jwk.Cache
}

func (a *AuthMiddleware) ValidateJWT(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, bearerSchema) {
		c.Header(wwwAuthenticateHeader, a.getAuthenticationUrl())
		responses.OCIUnauthorizedError(c)
		return
	}

	jwkSet, err := a.JwkCache.Lookup(c, a.Config.GetCognitoJWKUrl())
	if err != nil {
		responses.OCIInternalServerError(c)
		return
	}

	jwtToken := header[len(bearerSchema):]
	_, err = jwt.Parse([]byte(jwtToken), jwt.WithKeySet(jwkSet))
	if err != nil {
		c.Header(wwwAuthenticateHeader, a.getAuthenticationUrl())
		responses.OCIUnauthorizedError(c)
	}

	c.Next()
}

func (a *AuthMiddleware) getAuthenticationUrl() string {
	return fmt.Sprintf(`Bearer realm="https://%s/v2/login",service="registry.io"`, a.Config.BaseURL)
}
