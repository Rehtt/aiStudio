package redis

import (
	"aiStudio/internal/conf"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
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

var (
	lockTime = time.Minute * 5
	lockWait = time.Millisecond * 500
)

func Lock(ctx context.Context, key string, value any, exp ...time.Duration) error {
	t := lockTime
	if len(exp) != 0 {
		t = exp[0]
	}
	wait := time.NewTicker(lockWait)
	for {
		ok, err := DB.SetNX(ctx, key, value, t).Result()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		<-wait.C
	}
}
func UnLock(ctx context.Context, key string) error {
	return DB.Del(ctx, key).Err()
}
