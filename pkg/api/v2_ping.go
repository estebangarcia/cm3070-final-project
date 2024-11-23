package api

import "net/http"

type V2PingHandler struct{}

func (h *V2PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
