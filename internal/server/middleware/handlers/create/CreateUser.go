package handlers

import (
	"BankApp/internal/db"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type requestData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

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
func CreateNewUser(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, _ := c.MustGet("logger").(*zerolog.Logger)
		var newUserData requestData
		if err := c.ShouldBindJSON(&newUserData); err != nil {
			log.Error().Err(err).Msg("Failed to bind JSON for CreateUser")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		newUser := db.CreateUserParams{
			Email:    newUserData.Email,
			Password: newUserData.Password,
		}
		if _, err := mail.ParseAddress(newUser.Email); err != nil {
			log.Error().Err(err).Msg("Invalid email format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			c.Abort()
			return
		}
		id, err := queries.CreateUser(c.Request.Context(), newUser)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			c.Abort()
			return
		} else {
			log.Info().Dict("new_user_data", zerolog.Dict().
				Str("email", newUser.Email)).
				Msg("User in users table created successfully")
		}
		userProfile := db.CreateUserProfileParams{
			UsersID:   id,
			FirstName: newUserData.FirstName,
			LastName:  newUserData.LastName,
		}
		if _, err := queries.CreateUserProfile(c.Request.Context(), userProfile); err != nil {
			log.Error().Err(err).Msg("Failed to create user profile")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user profile"})
			c.Abort()
			return
		} else {
			log.Info().Dict("new_user_profile_data", zerolog.Dict().
				Int32("user_id", userProfile.UsersID)).
				Msg("User profile in user_profiles table created successfully")
		}
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
		log.Info().Msg("User creation process completed successfully")
	}
}
