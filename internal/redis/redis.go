package redis

import (
	"aiStudio/internal/conf"
	"context"
	"github.com/redis/go-redis/v9"
)

var DB *redis.Client

func New(c *conf.Redis) error {
	DB = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Username: c.Username,
		Password: c.Password,
		DB:       c.DB,
	})
	return DB.Ping(context.Background()).Err()
}
