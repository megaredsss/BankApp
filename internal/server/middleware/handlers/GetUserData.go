package handlers

import (
	"BankApp/internal/db"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func GetUsersBalanceHandler(c *gin.Context, queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, _ := c.MustGet("logger").(*zerolog.Logger)
		id := c.Param("id")
		if id == "" {
			log.Error().Msg("Id is empty")
			c.JSON(http.StatusOK, gin.H{"error": "ID is required"})
			c.Abort()
		}
		idInt, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			log.Error().Err(err).Msg("Invalid user ID format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			c.Abort()
		}
		balance, err := queries.GetUserBalance(c.Request.Context(), int32(idInt))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Error().Err(err).Msg("User not found")
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				c.Abort()
			} else {
				log.Error().Err(err).Msg("Failed to get user balance")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user balance"})
				c.Abort()
			}
		}
		c.JSON(http.StatusOK, gin.H{"balance": balance})
		log.Info().Int32("user_id", int32(idInt)).
			Msg("User balance retrieved successfully")
	}
}
