package handlers

import (
	"BankApp/db"
	"BankApp/resources/models"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := c.ShouldBindJSON(&inputUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validate.Struct(inputUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := mail.ParseAddress(inputUserData.Email); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := db.GetDB().Where("email = ?", inputUserData.Email).First(&inputUserData).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "User already exist"})
		return
	}
	newUser := models.UserDb{Email: inputUserData.Email, FirstName: inputUserData.FirstName, SecondName: inputUserData.SecondName, Balance: inputUserData.Balance, Password: inputUserData.Password}
	db.GetDB().Create(&newUser)
	c.JSON(http.StatusCreated, gin.H{"data": newUser})
}
