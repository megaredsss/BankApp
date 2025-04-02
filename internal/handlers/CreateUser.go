package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Creating User
func CreateUser(c *gin.Context) {
	var inputUserData models.UserDb
	if err := c.ShouldBindJSON(&inputUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser := models.UserDb{FirstName: inputUserData.FirstName, SecondName: inputUserData.SecondName, ThirdName: inputUserData.ThirdName, Balance: inputUserData.Balance, Password: inputUserData.Password}
	db.GetDB().Create(&newUser)
	c.JSON(http.StatusOK, gin.H{"data": newUser})
}
