package handlers

import (
	"BankApp/internal/db"
	"errors"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func CreateNewUser(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, _ := c.MustGet("logger").(*zerolog.Logger)
		userDataGet, exist := c.Get("userData")
		if !exist {
			log.Err(errors.New("value not found in context")).Msg("CreateNewUser data not found in gin context")
		}
		userData, ok := userDataGet.(RequestData)
		if !ok {
			log.Err(errors.New("error in type assertion")).Msg("Wrong type of userData in gin Context")
		}
		newUser := db.CreateUserParams{
			Email:    userData.Email,
			Password: userData.Password,
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
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
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
