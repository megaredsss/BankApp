package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"github.com/gin-gonic/gin"
)

func GetUserData(c *gin.Context) {
	var user []models.User
	db.GetDB().Find(&user)
	c.JSON(200, user)
}
