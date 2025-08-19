package handlers

import (
	"BankApp/internal/db"
	"BankApp/internal/server/middleware/logger"
	jwtPack "BankApp/pkg/jwt"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func ExpireSession(c *gin.Context, secret *jwtPack.SecretService) {
	log, ok := logger.GetLoggerFromContext(c)
	if !ok {
		fmt.Println("logger doesn't exist in context, func ExpireSession")
	}
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		log.Error().Err(err).Msg("No token in cookie")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}
	tokenStatus, err := secret.VerifyJWT(tokenString)
	if !tokenStatus {
		log.Error().Err(err).Msg("Failed to verify JWT")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
	}
	c.SetCookie("jwt", tokenString, -1, "/", "", false, true)
}

func LoginUser(c *gin.Context, queries *db.Queries, jwt *jwtPack.SecretService) {
	log, ok := logger.GetLoggerFromContext(c)
	if !ok {
		fmt.Println("logger doesn't exist in context, func LoginUser")
		return
	}
	var userData db.User
	if err := c.ShouldBindJSON(&userData); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON for LoginUser")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		c.Abort()
	}
	id, err := queries.GetUserIDByEmail(c.Request.Context(), userData.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user ID by email")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID by email"})
		c.Abort()
	}
	if id == 0 {
		log.Error().Msg("User not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
	}
	log.Info().Int32("user_id", id).Msg("User found")
	token, err := jwt.CreateJWT(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate JWT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		c.Abort()
	}
	log.Info().Dict("token_data", zerolog.Dict().
		Int32("user_id", id).
		Str("token", token)).
		Msg("JWT generated successfully")
	c.SetCookie("jwt", token, 3600, "/", "", false, true)
	log.Info().Msgf("User with id=%d logged in successfully", id)
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
	c.Next()
}
