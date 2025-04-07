package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TO-DO TEST THIS
func SendMoney(c *gin.Context) {
	var inputData models.TransferUser
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sender := models.UserDb{FirstName: inputData.SenderUser.FirstName, SecondName: inputData.SenderUser.SecondName, ThirdName: inputData.SenderUser.ThirdName}
	receiver := models.UserDb{FirstName: inputData.ReceiverUser.FirstName, SecondName: inputData.ReceiverUser.SecondName, ThirdName: inputData.ReceiverUser.ThirdName}
	amount := inputData.Amount
	if amount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
	}
	if err := db.GetDB().Where("first_name = ? AND second_name = ? AND third_name = ?", sender.FirstName, sender.SecondName, sender.ThirdName).First(&sender).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender not found"})
	}
	if err := db.GetDB().Where("first_name = ? AND second_name = ? AND third_name = ?", receiver.FirstName, receiver.SecondName, receiver.ThirdName).First(&receiver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
	}
	if int(sender.Balance) < amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender doesn't have enough money"})
	}
	db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&sender).Where("id = ? AND balance >= 0", sender.ID).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}
		if err := tx.Model(&receiver).Where("id = ? AND balance >= 0", receiver.ID).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}
		return nil
	})
}
