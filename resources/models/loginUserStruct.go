package models

type LoginUser struct {
	Email    string `binding:"required,email"`
	Password string `binding:"required,min=6"`
}
