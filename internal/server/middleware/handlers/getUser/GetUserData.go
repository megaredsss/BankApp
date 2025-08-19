package handlers

import (
	"BankApp/internal/db"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func GetUsersBalanceHandler(c *gin.Context, queries *db.Queries) {
	log, _ := c.MustGet("logger").(*zerolog.Logger)
	id := c.Param("id")
	if id == "" {
		log.Error().Msg("Id is empty")
		c.JSON(400, gin.H{"error": "ID is required"})
		c.Abort()
	}
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("Invalid user ID format")
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		c.Abort()
	}
	balance, err := queries.GetUserBalance(c.Request.Context(), int32(idInt))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("User not found")
			c.JSON(404, gin.H{"error": "User not found"})
			c.Abort()
		} else {
			log.Error().Err(err).Msg("Failed to get user balance")
			c.JSON(500, gin.H{"error": "Failed to get user balance"})
			c.Abort()
		}
	}
	c.JSON(200, gin.H{"balance": balance})
	log.Info().Int32("user_id", int32(idInt)).
		Msg("User balance retrieved successfully")
}
