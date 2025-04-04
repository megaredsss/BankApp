package redisPack

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

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
func GetRedis() *RedisClient { return client }
