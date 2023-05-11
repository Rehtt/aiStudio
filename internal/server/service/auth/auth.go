package auth

import (
	"aiStudio/internal/redis"
	model2 "aiStudio/internal/server/service/model"
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
	"strings"
)

const (
	ExternalAuthKey = "external:token:"
	ExternalLockKey = "external:lock:"
)

var whiteList = []string{
	"/progress/",
	"/image/",
}

func ExternalAuth() goweb.HandlerFunc {
	return func(ctx *goweb.Context) {
		var pass bool
		defer func() {
			if !pass {
				ctx.Stop()
			}
		}()
		var inWhiteList bool
		for _, w := range whiteList {
			if strings.Contains(ctx.Request.URL.Path, w) {
				inWhiteList = true
			}
		}

		token := ctx.Request.URL.Query().Get("token")
		if token == "" {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		key := ExternalAuthKey + token
		lockKey := ExternalLockKey + token

		number, _ := redis.DB.Get(ctx, key).Int()
		if !inWhiteList {
			if number <= 0 {
				ctx.WriteJSON(model2.CodeMap[model2.ResError], http.StatusBadRequest)
				return
			}

			// 加锁
			if err := redis.Lock(ctx, lockKey, "1"); err != nil {
				ctx.WriteJSON(model2.CodeMap[model2.ServerBad], http.StatusBadGateway)
				return
			}
		}

		number, _ = redis.DB.Get(ctx, key).Int()

		ctx.SetValue("info", &model2.ExternalInfo{
			Token:              token,
			RedisKey:           key,
			LockKey:            lockKey,
			Number:             number,
			ExpirationDuration: redis.DB.TTL(ctx, key).Val(),
		})
		pass = true
		return
	}
}
func Auth() goweb.HandlerFunc {
	return func(ctx *goweb.Context) {

	}
}
