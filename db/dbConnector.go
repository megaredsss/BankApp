package db

import (
	"BankApp/resources/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

// Get environment variables
// for the PostgresSQL
// return username, password, dbName
func getEnv() (int, string, string, string, string) {
	e := godotenv.Load(".env")
	if e != nil {
		fmt.Print(e)
	}
	host := os.Getenv("DB_HOST")
	port := 5432
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	return port, host, username, password, dbName
}

// Connect to PostgresSQL
// opening connection and do AutoMigration
func connectToDatabase(port int, host, username, password, dbName string) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", host, port, username, dbName, password)
	fmt.Println(dsn)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
	}

	database = conn
	err = database.Debug().AutoMigrate(&models.UserDb{})
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
