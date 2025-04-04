package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"net/http"

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
	newUser := models.UserDb{FirstName: inputUserData.FirstName, SecondName: inputUserData.SecondName, ThirdName: inputUserData.ThirdName, Balance: inputUserData.Balance, Password: inputUserData.Password}
	db.GetDB().Create(&newUser)
	c.JSON(http.StatusOK, gin.H{"data": newUser})
}
