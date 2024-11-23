package responses

import (
	"encoding/json"
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
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(OCIErrorResponse{
		Errors: []OCIErrorDetail{
			{
				Code:    "UNAUTHORIZED",
				Message: "Authentication is required",
				Detail:  nil,
			},
		},
	})
}

func OCIInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(OCIErrorResponse{
		Errors: []OCIErrorDetail{
			{
				Code:    "INTERNAL_ERROR",
				Message: "Internal Server Error",
				Detail:  nil,
			},
		},
	})
}
