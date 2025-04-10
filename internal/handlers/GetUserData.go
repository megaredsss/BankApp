package handlers

import (
	"BankApp/db"
	"BankApp/pkg/redisPack"
	"BankApp/resources/models"
	"context"
	"errors"
	"net/mail"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get user's Balance
// by first name, second name and third name
func GetUsersBalance(c *gin.Context) {
	var user models.UserDb
	usersEmail := c.DefaultQuery("email", "")
	if _, err := mail.ParseAddress(usersEmail); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	usersID, ok := redisPack.GetRedis().Get(context.Background(), usersEmail).Result()
	if ok != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}
	if err := db.GetDB().Select("balance").Where("id = ?", usersID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(400, gin.H{"error": gorm.ErrRecordNotFound})
		} else {
			c.JSON(400, gin.H{"error": err})
		}
	} else {
		c.JSON(200, gin.H{"User found": user.FirstName + user.SecondName})
	}
	c.JSON(200, user.Balance)
}
