package api

import "github.com/gin-gonic/gin"

type V2LoginHandler struct{}

func (h *V2LoginHandler) Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"access_token": "asd123",
	})
}
