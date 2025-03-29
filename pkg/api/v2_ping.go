package api

import "net/http"

type V2PingHandler struct{}

// This handler just returns a 200 OK and is an endpoint used by docker and
// oci tooling to ping the registry to check if the credentials are valid and
// if the registry is up and available
func (h *V2PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
