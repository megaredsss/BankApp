package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func DeleteUser(c *gin.Context) {
	var user models.UserDb
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	if err := db.GetDB().Where("email = ? AND password = ?", user.Email, user.Password).Delete(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
