package handlers

import (
	"BankApp/internal/db"
	jwtPack "BankApp/pkg/jwt"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func DeleteUser(queries *db.Queries, secret *jwtPack.SecretService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, _ := c.MustGet("logger").(*zerolog.Logger)
		token, err := c.Cookie("jwt")
		if err != nil {
			log.Error().Err(err).Msg("jwt token doesn't exist")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "bad token"})
			c.Abort()
		}
		id, err := secret.GetIdFromClaims(token)
		if err != nil {
			log.Error().Err(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bad token"})
			c.Abort()
		}
		if err := queries.DeleteUser(c, id); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Error().Err(err).Msg("User not found")
				c.JSON(404, gin.H{"error": "User not found"})
				c.Abort()
			} else {
				log.Error().Err(err).Int32("Error during delete operation for user with id=", id)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during delete operation"})
				c.Abort()
			}
		}
		log.Info().Msgf("User with id=%d deleted successfully", id)
		c.JSON(200, gin.H{"message": "User deleted successfully"})
	}
}
