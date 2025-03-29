package main

import (
	"BankApp/internal/handlers"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()

	router.GET("/UserData", handlers.GetUserData)
	router.POST("/CreateUser", handlers.CreateUser)
	// @Summary Add a new pet to the store
	// @Description get string by ID
	// @ID get-string-by-int
	// @Accept  json
	// @Produce  json
	// @Param   some_id     path    int     true        "Some ID"
	// @Success 200 {string} string  "ok"
	// @Router /string/{some_id} [get]
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Error on Run")
	}

}
