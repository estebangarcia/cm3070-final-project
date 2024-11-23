package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
	"github.com/gin-gonic/gin"
)

type V2LoginHandler struct {
	Config        *config.AppConfig
	CognitoClient *cognitoidentityprovider.Client
}

func (h *V2LoginHandler) Login(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		responses.OCIUnauthorizedError(c)
		return
	}

	uV, _ := c.Get("username")
	pV, _ := c.Get("password")

	username := uV.(string)
	password := pV.(string)

	mac := hmac.New(sha256.New, []byte(h.Config.Cognito.ClientSecret))
	mac.Write([]byte(username + h.Config.Cognito.ClientId))

	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	params := map[string]string{
		"USERNAME":    username,
		"PASSWORD":    password,
		"SECRET_HASH": secretHash,
	}

	// Initiate authentication request
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserPasswordAuth,
		ClientId:       aws.String(h.Config.Cognito.ClientId),
		AuthParameters: params,
	}

	// Call Cognito
	output, err := h.CognitoClient.InitiateAuth(c, input)
	if err != nil {
		log.Println(err.Error())
		responses.OCIUnauthorizedError(c)
		return
	}

	c.JSON(200, gin.H{
		"token":      output.AuthenticationResult.AccessToken,
		"expires_in": output.AuthenticationResult.ExpiresIn,
	})
}
