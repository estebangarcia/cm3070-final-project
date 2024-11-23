package responses

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int32  `json:"expires_in"`
}
