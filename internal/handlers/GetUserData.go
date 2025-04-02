package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get user's Balance
// by first name, second name and third name
func GetUserData(c *gin.Context) {
	var user models.UserDb
	usersFirstName := c.DefaultQuery("first_name", "")
	usersSecondName := c.DefaultQuery("second_name", "")
	usersThirdName := c.DefaultQuery("third_name", "")
	if err := db.GetDB().Select("balance").Where("first_name = ? AND second_name = ? AND third_name = ?", usersFirstName, usersSecondName, usersThirdName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Not found")
		} else {
			fmt.Println("Bad query", err)
		}
	} else {
		fmt.Println("User found", user.FirstName, user.SecondName, user.ThirdName)
	}
	c.JSON(200, user.Balance)
}
