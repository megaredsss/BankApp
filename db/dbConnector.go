package db

import (
	"BankApp/resources/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var database *gorm.DB

// Get environment variables
// for the PostgresSQL
// return username, password, dbName
func getEnv() (string, string, string) {
	e := godotenv.Load(".env")
	if e != nil {
		fmt.Print(e)
	}
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	return username, password, dbName
}

// Connect to PostgresSQL
// opening connection and do AutoMigration
func connectToDatabase(username, password, dbName string) {
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

func init() {
	connectToDatabase(getEnv())
}

func GetDB() *gorm.DB {
	return database
}
