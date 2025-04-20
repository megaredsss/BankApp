package models

// Struct for creating user
type User struct {
	FirstName  string `binding:"required,min=2"`
	SecondName string `binding:"required,min=2"`
	ThirdName  string `binding:"required,min=2"`
	Balance    uint   `binding:"number"`
}
