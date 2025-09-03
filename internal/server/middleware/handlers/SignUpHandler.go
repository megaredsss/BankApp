package handlers

import (
	"BankApp/internal/db"
	"BankApp/internal/server/middleware/logger"
	jwtPack "BankApp/pkg/jwt"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type RequestData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func SignUpUser(queries *db.Queries, secret *jwtPack.SecretService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, ok := logger.GetLoggerFromContext(c)
		if !ok {
			fmt.Println("logger doesn't exist in context, func LoginUser")
			return
		}
		var userDataFromRequest RequestData
		if err := c.ShouldBindJSON(&userDataFromRequest); err != nil {
			log.Error().Err(err).Msg("Failed to bind JSON for SignUpUser")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			c.Abort()
			return
		}
		userData := db.CreateUserParams{
			Email: userDataFromRequest.Email,
		}
		id, err := queries.GetUserIDByEmail(c.Request.Context(), userData.Email)
		if err != pgx.ErrNoRows {
			log.Error().Err(err).Msg("Failed to get user ID by email")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from Database"})
			c.Abort()
			return
		}
		if id != 0 {
			log.Error().Msg("User found but should't")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User already exists"})
			c.Abort()
			return
		}
		if token, err := c.Cookie("jwt"); err == nil && token != "" {
			log.Error().Err(err).Msg("jwt token already exist")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you are already logged"})
			c.Abort()
		}
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
		c.Set("userData", userDataFromRequest)
		log.Info().Msgf("User with id=%d logged in successfully", id)
		c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
		c.Next()
	}

}
