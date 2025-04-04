package redisPack

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	Set(ctx context.Context, key string, value interface{}, expTime time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Ping(ctx context.Context) *redis.StatusCmd
}
type RedisClient struct {
	client *redis.Client
}

var client *RedisClient

func (rds *RedisClient) Set(ctx context.Context, key string, value interface{}, expTime time.Duration) *redis.StatusCmd {
	return rds.client.Set(ctx, key, value, expTime)
}

func (rds *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return rds.client.Get(ctx, key)
}

func (rds *RedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	return rds.client.Ping(ctx)
}
