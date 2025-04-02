package handlers

import (
	"BankApp/db"
	jwtPack "BankApp/jwt"
	"BankApp/resources/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginHandler(c *gin.Context) {
	var userData models.LoginUser
	var userDb models.UserDb
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	if err := db.GetDB().Where("first_name = ? AND second_name = ? AND third_name = ?", userData.FirstName, userData.SecondName, userData.ThirdName).First(&userDb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Not found")
		} else {
			fmt.Println("Bad query", err)
		}
	} else {
		fmt.Println("User found", userData.FirstName, userData.SecondName, userData.ThirdName)
	}
}
func tokenChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token"})
			c.Abort()
			return
		}
		tokenStatus, err := jwtPack.VerifyJWT(tokenString)
		if !tokenStatus {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Next()
	}
}

func authChecker(c *gin.Context) {

}
