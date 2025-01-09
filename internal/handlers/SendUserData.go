package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Creating User
func CreateUser(c *gin.Context) {
	var inputUserData models.CreateUserStruct
	if err := c.ShouldBindJSON(&inputUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser := models.User{FirstName: inputUserData.FirstName, SecondName: inputUserData.SecondName, ThirdName: inputUserData.ThirdName, Balance: inputUserData.Balance}
	db.GetDB().Create(&newUser)
	c.JSON(http.StatusOK, gin.H{"data": newUser})
}
