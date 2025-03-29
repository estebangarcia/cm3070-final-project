package api

import (
	"encoding/json"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

type HealthHandler struct {
}

// Dummy handler to be used for health checks
func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.HealthResponse{
		Status: "ok",
	})
}
