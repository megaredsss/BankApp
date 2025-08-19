package main

import (
	"BankApp/internal/config"
	"BankApp/internal/db"
	handlers "BankApp/internal/server/middleware/handlers/create"
	"BankApp/internal/server/middleware/logger"
	"BankApp/internal/server/middleware/requestId"
	jwtPack "BankApp/pkg/jwt"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupLogger(env string) {
	switch env {
	case "dev":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	case "prod":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {
	cfg := config.Loader()
	fmt.Println("Config loaded successfully:", cfg)
	setupLogger(cfg.Env)
	conn, err := db.ConnectToDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
		return
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()
	log.Info().Msgf("Connected to database %s on %s:%d", cfg.Database.Name, cfg.Database.Host, cfg.Database.Port)
	queries := conn.CreateNew()
	jwt := jwtPack.NewSecretService(cfg.SecretKey)
	log.Info().Msg("JWT service initialized")

	router := gin.New()
	router.Use(requestId.GenerateRequestID())
	router.Use(gin.Recovery())
	router.Use(logger.LoggerMiddleware(&log.Logger))

	_ = queries
	_ = jwt
	router.POST("/CreateUser", handlers.CreateNewUser(queries))
	// router.POST("/Login", handlers.LoginHandler)
	// router.GET("/GetBalance", handlers.TokenChecker, handlers.GetUsersBalance)
	// router.DELETE("/DeleteUser", handlers.DeleteUser)
	// router.PUT("/SendMoney", handlers.TokenChecker, handlers.SendMoney)
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	err = router.Run(":8080")
	if err != nil {
		fmt.Println("Error on Run")
	}

}
