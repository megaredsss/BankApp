package redisPack

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// @interface RedisInterface - сам интерфейс для работы с Redis
// @method Set - метод для установки значения по ключу
// @method Get - метод для получения значения по ключу
// @method Ping - метод для проверки соединения с Redis
type RedisInterface interface {
	Set(ctx context.Context, key string, value interface{}, expTime time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, key string) *redis.IntCmd
	Ping(ctx context.Context) *redis.StatusCmd
}

// @struct RedisClient - структура для работы с Redis
type RedisClient struct {
	client *redis.Client
}

// client - клиент Redis
var client *RedisClient

// @method Set - метод для установки значения по ключу
// @param ctx - контекст, в котором выполняется операция
// @param key - ключ, по которому будет установлено значение
// @param value - значение, которое будет установлено по ключу
// @param expTime - время жизни ключа
// @return *redis.StatusCmd - результат выполнения операции
func (rds *RedisClient) Set(ctx context.Context, key string, value interface{}, expTime time.Duration) *redis.StatusCmd {
	return rds.client.Set(ctx, key, value, expTime)
}

// @method Get - метод для получения значения по ключу
// @param ctx - контекст, в котором выполняется операция
// @param key - ключ, по которому будет получено значение
// @return *redis.StringCmd - результат выполнения операции
func (rds *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return rds.client.Get(ctx, key)
}

func (rds *RedisClient) Del(ctx context.Context, key string) *redis.IntCmd {
	return rds.client.Del(ctx, key)
}

// @method Ping - метод для проверки соединения с Redis
// @param ctx - контекст, в котором выполняется операция
// @param key - ключ, по которому будет получено значение
// @return *redis.StatusCmd - результат выполнения операции
func (rds *RedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	return rds.client.Ping(ctx)
}
