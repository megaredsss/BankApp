package db

import (
	"BankApp/resources/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var database *gorm.DB

func init() {

	e := godotenv.Load("../.env")
	if e != nil {
		fmt.Print(e)
	}
	username := os.Getenv("db_user")
	password := os.Getenv("db_password")
	dbName := os.Getenv("db_name")
	dsn := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s", username, dbName, password)
	fmt.Println(dsn)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
	}

	database = conn
	err = database.Debug().AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("AutoMigrate error")
	}
}

func GetDB() *gorm.DB {
	return database
}
