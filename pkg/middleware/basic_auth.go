package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
	"github.com/gin-gonic/gin"
)

const basicSchema = "Basic "

type ExtractBasicCredentialsMiddleware struct {
}

func (a *ExtractBasicCredentialsMiddleware) Validate(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, basicSchema) {
		responses.OCIUnauthorizedError(c)
		return
	}

	basicCreds := header[len(basicSchema):]

	decoded, err := base64.StdEncoding.DecodeString(basicCreds)
	if err != nil {
		responses.OCIUnauthorizedError(c)
		return
	}

	usernamePassword := strings.Split(string(decoded), ":")
	if len(usernamePassword) < 2 {
		responses.OCIUnauthorizedError(c)
		return
	}

	username := usernamePassword[0]
	password := usernamePassword[1]

	if username == "" || password == "" {
		responses.OCIUnauthorizedError(c)
		return
	}

	c.Set("username", username)
	c.Set("password", password)
	c.Next()
}
