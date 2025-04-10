package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"net/mail"

	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	var user models.UserDb
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	result := db.GetDB().Where("email = ? AND password = ?", user.Email, user.Password).Delete(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
