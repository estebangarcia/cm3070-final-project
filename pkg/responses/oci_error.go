package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OCIErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
}

type OCIErrorResponse struct {
	Errors []OCIErrorDetail `json:"errors"`
}

func GenericOCIError(w http.ResponseWriter, code string, status int, message string, detail interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(OCIErrorResponse{
		Errors: []OCIErrorDetail{
			{
				Code:    code,
				Message: message,
				Detail:  detail,
			},
		},
	})
}

func OCIUnauthorizedError(w http.ResponseWriter) {
	GenericOCIError(w, "UNAUTHORIZED", http.StatusUnauthorized, "Authentication is required", nil)
}

func OCIInternalServerError(w http.ResponseWriter) {
	GenericOCIError(w, "INTERNAL_ERROR", http.StatusInternalServerError, "Internal Server Error", nil)
}

func OCIManifestUnknown(w http.ResponseWriter, reference string) {
	GenericOCIError(w, "MANIFEST_UNKNOWN", http.StatusNotFound, fmt.Sprintf("Manifest with reference '%s' not found", reference), nil)
}

func OCIBlobUnknown(w http.ResponseWriter, digest string) {
	GenericOCIError(w, "BLOB_UNKNOWN", http.StatusNotFound, fmt.Sprintf("Blob with digest '%s' not found", digest), nil)
}
