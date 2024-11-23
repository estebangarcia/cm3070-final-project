package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OCIErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
}

type OCIErrorResponse struct {
	Errors []OCIErrorDetail `json:"errors"`
}

func GenericOCIError(c *gin.Context, code string, status int, message string, detail interface{}) {
	c.JSON(status, OCIErrorResponse{
		Errors: []OCIErrorDetail{
			{
				Code:    code,
				Message: message,
				Detail:  detail,
			},
		},
	})
	c.Abort()
}

func OCIUnauthorizedError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, OCIErrorResponse{
		Errors: []OCIErrorDetail{
			{
				Code:    "UNAUTHORIZED",
				Message: "Authentication is required",
				Detail:  nil,
			},
		},
	})
	c.Abort()
}

func OCIInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, OCIErrorResponse{
		Errors: []OCIErrorDetail{
			{
				Code:    "INTERNAL_ERROR",
				Message: "Internal Server Error",
				Detail:  nil,
			},
		},
	})
	c.Abort()
}
