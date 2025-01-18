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

func OCIManifestBlobUnknown(w http.ResponseWriter, reference string) {
	GenericOCIError(w, "MANIFEST_BLOB_UNKNOWN", http.StatusNotFound, fmt.Sprintf("Manifest references non-existant digest '%s'", reference), nil)
}

func OCIBlobUnknown(w http.ResponseWriter, digest string) {
	GenericOCIError(w, "BLOB_UNKNOWN", http.StatusNotFound, fmt.Sprintf("Blob with digest '%s' not found", digest), nil)
}

func OCIBlobUploadInvalid(w http.ResponseWriter) {
	GenericOCIError(w, "BLOB_UPLOAD_INVALID", http.StatusBadRequest, "Failed to upload blob", nil)
}

func OCIBlobUploadUnknown(w http.ResponseWriter) {
	GenericOCIError(w, "BLOB_UPLOAD_UNKNOWN", http.StatusNotFound, "Blob upload session not found", nil)
}

func OCIRepositoryUnknown(w http.ResponseWriter, repositoryName string, asBadRequest bool) {
	var statusCode = http.StatusNotFound
	if asBadRequest {
		statusCode = http.StatusBadRequest
	}

	GenericOCIError(w, "REPOSITORY_UNKNOWN", statusCode, fmt.Sprintf("Repository with name '%s' not found", repositoryName), nil)
}

func OCIUnprocessableEntity(w http.ResponseWriter, message string) {
	GenericOCIError(w, "UNPROCESSABLE_ENTITY", http.StatusUnprocessableEntity, message, nil)
}

func OCITagUnknown(w http.ResponseWriter, repositoryName string, tagName string) {
	GenericOCIError(w, "TAG_UNKNOWN", http.StatusNotFound, fmt.Sprintf("Repository with name '%s' has no tag with name '%s'", repositoryName, tagName), nil)
}
