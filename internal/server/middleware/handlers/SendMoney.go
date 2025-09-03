package handlers

import (
	"BankApp/internal/db"
	"BankApp/internal/server/middleware/logger"
	jwtPack "BankApp/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type SendMoneyRequest struct {
	Amount     float64 `json:"amount" binding:"required"`
	ReceiverId int32   `json:"receiver" binding:"required"`
}

func numericToFloat64Converter(inputNumeric pgtype.Numeric) (float64, error) {
	var result float64
	if inputNumeric.Valid && !inputNumeric.NaN && inputNumeric.InfinityModifier == pgtype.Finite {
		inputNumericFloat8, err := inputNumeric.Float64Value()
		if err != nil {
			return 0, err
		}
		result = inputNumericFloat8.Float64
	}
	return result, nil
}
func checkBalanceForTransaction(c *gin.Context, queries *db.Queries, database *db.Database, senderId int32, amount float64) (bool, error) {
	var senderBalanceFloat float64
	senderBalance, err := queries.GetUserBalance(c, senderId)
	if err != nil {
		return false, err
	}
	senderBalanceFloat, err = numericToFloat64Converter(senderBalance)
	if err != nil {
		return false, err
	}
	if senderBalanceFloat < amount {
		return false, nil
	}
	return true, nil
}
func SendMoney(queries *db.Queries, database *db.Database, secret *jwtPack.SecretService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, ok := logger.GetLoggerFromContext(c)
		if !ok {
			log.Error().Msg("logger doesn't exist in context, func SendMoney")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		var request SendMoneyRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Error().Err(err).Msg("Invalid request data")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			c.Abort()
			return
		}
		var amount pgtype.Numeric
		if err := amount.Scan(request.Amount); err != nil {
			log.Error().Err(err).Msg("Failed to scan amount")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
			c.Abort()
			return
		}
		token, err := c.Cookie("jwt")
		if err != nil {
			log.Error().Err(err).Msg("No token in cookie")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token"})
			c.Abort()
			return
		}
		tokenStatus, err := secret.VerifyJWT(token)
		if err != nil || !tokenStatus {
			log.Error().Err(err).Msg("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		senderId, err := secret.GetIdFromClaims(token)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user ID from token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if balanceStatus, err := checkBalanceForTransaction(c, queries, database, senderId, request.Amount); err != nil || !balanceStatus {
			if err != nil {
				log.Error().Err(err).Msg("Failed to check balance")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			if !balanceStatus {
				log.Error().Msg("Not enough balance")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough balance"})
			}
			c.Abort()
			return
		}
		tx, err := database.Db.Begin(c)
		if err != nil {
			c.Abort()
			return
		}
		defer tx.Rollback(c)
		qtx := queries.WithTx(tx)
		if err := qtx.MoneyTransactionIncreaseByUserId(c, db.MoneyTransactionIncreaseByUserIdParams{
			ID:      request.ReceiverId,
			Balance: amount,
		}); err != nil {
			c.Abort()
			return
		}
		if err := qtx.MoneyTransactionDecreaseByUserId(c, db.MoneyTransactionDecreaseByUserIdParams{
			ID:      senderId,
			Balance: amount,
		}); err != nil {
			c.Abort()
			return
		}
		if err := tx.Commit(c); err != nil {
			log.Error().Err(err).Msg("Failed to commit transaction")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Money sent successfully"})
		log.Info().Int32("sender_id", senderId).
			Int32("receiver_id", request.ReceiverId).
			Float64("amount", request.Amount).
			Msg("Money sent successfully")
		c.Next()
	}
}
