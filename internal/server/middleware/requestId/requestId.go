package requestId

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

func GenerateRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader("X-Request-ID")
		if requestId == "" {
			bytes := make([]byte, 16)
			rand.Read(bytes)
			requestId = base64.StdEncoding.EncodeToString(bytes)
		}
		c.Set("request_id", requestId)
		c.Next()
	}
}
