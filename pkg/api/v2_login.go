package api

import (
	"encoding/json"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

type V2LoginHandler struct {
}

func (h *V2LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(string)
	expiresIn := r.Context().Value("expires_in").(int32)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.TokenResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	})
}
