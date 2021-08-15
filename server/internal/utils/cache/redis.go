package cache

import (
	"context"
	"errors"
	"github.com/carol-caires/udp-chat/configs"
	"github.com/go-redis/redis/v8"
)

type RedisImpl struct {
	conn *redis.Client
}

func NewRedisConn() (*RedisImpl, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.GetRedisAddr(),
		Password: "", // no password set // todo: set password
		DB:       0,  // use default DB
	})
	if rdb == nil {
		return nil, errors.New("failed to create a connection on Redis")
	}
	return &RedisImpl{conn: rdb}, nil
}

func (c *RedisImpl) Set(ctx context.Context, key, value string) (err error) {
	result := c.conn.Set(ctx, key, value, 0)
	return result.Err()
}
func (c *RedisImpl) Get(ctx context.Context, key string) (value string, err error) {
	result := c.conn.Get(ctx, key)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Result()
}