package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary      Создание нового пользователя
// @Description  Создание нового пользователя с заданными данными
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body models.UserDb true "Данные нового пользователя"
// @Success      200 {object} models.UserDb "Успешное создание пользователя"
// @Failure      400 {object} gin.H{"error": "Ошибка валидации данных"}
// @Failure      500 {object} gin.H{"error": "Ошибка сервера"}
// @Router       /CreateUser [post]
func CreateUser(c *gin.Context) {
	var inputUserData models.UserDb
	if err := c.ShouldBindJSON(&inputUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := mail.ParseAddress(inputUserData.Email); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	newUser := models.UserDb{Email: inputUserData.Email, FirstName: inputUserData.FirstName, SecondName: inputUserData.SecondName, Balance: inputUserData.Balance, Password: inputUserData.Password}
	db.GetDB().Create(&newUser)
	c.JSON(http.StatusOK, gin.H{"data": newUser})
}
