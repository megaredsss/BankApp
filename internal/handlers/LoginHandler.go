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
	fmt.Println(userData)
	if err := db.GetDB().Where("first_name = ? AND second_name = ? AND third_name = ?", userData.FirstName, userData.SecondName, userData.ThirdName).First(&userDb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Not found")
			return
		} else {
			fmt.Println("Bad query", err)
			return
		}
	} else {
		fmt.Println("User found", userData.FirstName, userData.SecondName, userData.ThirdName)
	}
	usersToken, err := jwtPack.CreateJWT(int(userDb.ID), userData.FirstName, userData.SecondName)
	if err != nil {
		fmt.Println("Error in creating JWT during login", err)
		return
	}
	fmt.Println(userData, usersToken)
	c.SetCookie("jwt", usersToken, 3600, "/", "", false, true)
}
func TokenChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")
		if err != nil {
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
