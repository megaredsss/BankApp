package handlers

import (
	"BankApp/db"
	jwtPack "BankApp/jwt"
	"BankApp/resources/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func loginHandler(c *gin.Context) {
	var userData models.LoginUser
	var userDb models.UserDb
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	if db.GetDB().Where("first_name = ? AND second_name = ? AND third_name = ?", userData.FirstName, userData.SecondName, userData.ThirdName).First(&userDb) != nil {

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
