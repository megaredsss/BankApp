package middleware

import (
	"BankApp/internal/server/middleware/logger"
	jwtPack "BankApp/pkg/jwt"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TokenChecker(secret *jwtPack.SecretService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, ok := logger.GetLoggerFromContext(c)
		if !ok {
			fmt.Println("logger doesn't exist in context, func TokenChecker")
		}
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			log.Error().Err(err).Msg("No token in cookie")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
		tokenStatus, err := secret.VerifyJWT(tokenString)
		if err != nil {
			log.Error().Err(err).Msg("Failed to verify JWT")
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
		}
		if !tokenStatus {
			log.Error().Msg(("Invalid token"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
		c.Next()
	}
}
