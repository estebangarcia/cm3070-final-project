package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

type V2PingHandler struct{}

func (h *V2PingHandler) Ping(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	log.Println(auth)
	c.AbortWithStatus(200)
}
