package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Get user's Balance
// by first name, second name and third name
func GetUserData(c *gin.Context) {
	var user models.User
	usersFirstName := c.DefaultQuery("first_name", "")
	usersSecondName := c.DefaultQuery("second_name", "")
	usersThirdName := c.DefaultQuery("third_name", "")
	err := db.GetDB().Select("balance").Where("first_name = ? AND second_name = ? AND third_name = ?", usersFirstName, usersSecondName, usersThirdName).First(&user).Error
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, user.Balance)
}
