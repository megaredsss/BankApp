package handlers

import (
	"BankApp/db"
	jwtPack "BankApp/jwt"
	"BankApp/resources/models"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	ctx := context.Background()
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	if err := validate.Struct(userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.GetDB().Where("email = ? AND password = ?", userData.Email, userData.Password).First(&userDb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusFound, gin.H{"message": "User found"})
	}
	usersToken, err := jwtPack.CreateJWT(int(userDb.ID), userData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(userData, usersToken)
	if err := jwtPack.SaveJWTInRedis(ctx, usersToken, userDb.ID, time.Hour); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	tokenStatus, err := jwtPack.VerifyJWT(tokenString)
	if !tokenStatus {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.Next()
}

func ExpireSession(c *gin.Context) {
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	tokenStatus, err := jwtPack.VerifyJWT(tokenString)
	if !tokenStatus {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("jwt", tokenString, -1, "/", "", false, true)
}
