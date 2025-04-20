package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SendMoney(c *gin.Context) {
	var inputData models.TransferUser
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in bind JSON"})
		return
	}
	sender := models.UserDb{Email: inputData.SenderUser.Email}
	receiver := models.UserDb{Email: inputData.ReceiverUser.Email}
	if err := db.GetDB().Where("email = ?", sender.Email).First(&sender).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender doesn't exist"})
		return
	}
	if err := db.GetDB().Where("email = ?", receiver.Email).First(&receiver).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Receiver doesn't exist"})
		return
	}
	amount := inputData.Amount
	if amount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}
	if err := db.GetDB().Where("email = ?", sender.Email).First(&sender).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.GetDB().Where("email = ?", receiver.Email).First(&receiver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	senderBalance, _ := strconv.ParseFloat(sender.Balance, 64)
	if senderBalance < amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender doesn't have enough money"})
		return
	}
	if err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&sender).Where("id = ? AND balance >= ?", sender.ID, amount).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}
		fmt.Println(1)
		if err := tx.Model(&receiver).Where("id = ?", receiver.ID).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}
