package redisPack

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// getEnv() получает значения необходимые для подключения к Redis
// @param string, string, int - адрес, пароль, номер базы данных
// @return string, string, int - адрес, пароль, номер базы данных
func getEnv() (string, string, int) {
	e := godotenv.Load(".env")
	if e != nil {
		fmt.Print(e)
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println(err)
	}
	return redisAddr, redisPassword, redisDb
}

// ConnectRedis() подключается к Redis
// @param string, string, int - адрес, пароль, номер базы данных
// @return error - ошибка подключения к Redis
func ConnectRedis(Addr string, Password string, DB int) error {
	rds := redis.NewClient(
		&redis.Options{
			Addr:     Addr,
			Password: Password,
			DB:       DB,
		},
	)
	err := rds.Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	client = &RedisClient{client: rds}
	return nil
}

func init() {
	ConnectRedis(getEnv())
}

// GetRedis() возвращает клиент Redis
// @return *RedisClient - клиент Redis
func GetRedis() *RedisClient { return client }
