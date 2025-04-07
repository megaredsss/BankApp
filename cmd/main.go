package main

import (
	"BankApp/internal/handlers"
	jwtPack "BankApp/jwt"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	jwtPack.InitSecretKey()
	router := gin.Default()
	router.GET("/UserData", handlers.GetUserData)
	router.POST("/CreateUser", handlers.CreateUser)
	router.POST("/Login", handlers.LoginHandler)
	router.GET("/Balance", handlers.TokenChecker, handlers.GetBalance)
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Error on Run")
	}

}
