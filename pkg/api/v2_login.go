package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

type V2LoginHandler struct {
	Config        *config.AppConfig
	CognitoClient *cognitoidentityprovider.Client
}

func (h *V2LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		responses.OCIUnauthorizedError(w)
		return
	}

	uV := r.Context().Value("username")
	pV := r.Context().Value("password")

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
	output, err := h.CognitoClient.InitiateAuth(r.Context(), input)
	if err != nil {
		log.Println(err.Error())
		responses.OCIUnauthorizedError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.TokenResponse{
		Token:     *output.AuthenticationResult.AccessToken,
		ExpiresIn: output.AuthenticationResult.ExpiresIn,
	})
}
