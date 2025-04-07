package handlers

import (
	"BankApp/db"
	jwtPack "BankApp/jwt"
	"BankApp/pkg/redisPack"
	"BankApp/resources/models"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginHandler godoc
// @Summary      Аутентификация пользователя
// @Description  Обрабатывает запросы на аутентификацию, проверяя наличие пользователя в базе данных и создавая JWT-токен и сохраняя его в Redis.
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        userData body models.LoginUser true "Данные пользователя для аутентификации"
// @Success      200 {object} gin.H{"message": "Аутентификация успешна"} "Успешная аутентификация"
// @Failure      400 {object} gin.H{"error": "Неверный запрос"} "Ошибка в запросе"
// @Failure      401 {object} gin.H{"error": "Неверное имя пользователя"} "Пользователь не найден"
// @Failure      500 {object} gin.H{"error": "Ошибка сервера"} "Внутренняя ошибка сервера"
// @Router       /Login [post]
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
	redisPack.GetRedis().Set(c, usersToken, userDb.ID, time.Hour)
	c.SetCookie("jwt", usersToken, 3600, "/", "", false, true)
}

// TokenChecker godoc
// @Summary      Проверка JWT-токена
// @Description  Middleware для проверки наличия и валидности JWT-токена.
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Failure      401 {object} gin.H{"error": "No token"} "Отсутствует токен"
// @Failure      401 {object} gin.H{"error": "Invalid token"} "Неверный токен"
// @Router       / [any]
func TokenChecker(c *gin.Context) {
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token in Cookie"})
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
