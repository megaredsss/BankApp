package handlers

import (
	"BankApp/internal/db"
	"BankApp/internal/server/middleware/logger"
	jwtPack "BankApp/pkg/jwt"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type signInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

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
		return
	}
	tokenStatus, err := secret.VerifyJWT(tokenString)
	if !tokenStatus {
		log.Error().Err(err).Msg("Failed to verify JWT")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	c.SetCookie("jwt", tokenString, -1, "/", "", false, true)
}

func SignInUser(queries *db.Queries, secret *jwtPack.SecretService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, ok := logger.GetLoggerFromContext(c)
		if !ok {
			fmt.Println("logger doesn't exist in context, func LoginUser")
			return
		}
		if _, exist := c.Get("jwt"); exist {
			log.Info().Msg("User already logged")
			c.JSON(http.StatusContinue, gin.H{"info": " you are already logged"})
			c.Next()
			return
		}

		var signInRequestData signInRequest
		if err := c.ShouldBindQuery(&signInRequestData); err != nil {
			log.Error().Err(err).Msg("Failed to bind sign in data")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			c.Abort()
			return
		}
		log.Info().Dict("signInData", zerolog.Dict().
			Str("email", signInRequestData.Email)).
			Msg("Sign in data bound successfully")

		userData := db.User{
			Email:    signInRequestData.Email,
			Password: signInRequestData.Password,
		}
		if userData.Email == "" || userData.Password == "" {
			log.Error().Err(errors.New("empty data")).Msg("email or password are empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "email or password are empty"})
			c.Abort()
			return
		}
		id, err := queries.GetUserIDByEmail(c.Request.Context(), userData.Email)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user ID by email")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID by email"})
			c.Abort()
			return
		}
		if id == 0 {
			log.Error().Msg("User not found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		log.Info().Int32("user_id", id).Msg("User found")
		token, err := secret.CreateJWT(id)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate JWT")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
			c.Abort()
			return
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
}
