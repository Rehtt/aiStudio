package auth

import (
	"aiStudio/internal/redis"
	"aiStudio/internal/service/model"
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
)

const (
	ExternalAuthKey = "external-token::"
	ExternalLockKey = "external-lock::"
)

func ExternalAuth() goweb.HandlerFunc {
	return func(ctx *goweb.Context) {
		var pass bool
		defer func() {
			if !pass {
				ctx.Stop()
			}
		}()

		token := ctx.Request.URL.Query().Get("token")
		if token == "" {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		key := ExternalAuthKey + token
		lockKey := ExternalLockKey + token
		redis.Lock(ctx, lockKey, 1)
		number, err := redis.DB.Get(ctx, key).Int()
		if err != nil {
			ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
			return
		}
		if number > 0 {
			ctx.SetValue("info", &model.ExternalInfo{
				Key:                key,
				Number:             number,
				ExpirationDuration: redis.DB.TTL(ctx, key).Val(),
				LockKey:            lockKey,
			})
			pass = true
			return
		}
		ctx.WriteJSON(model.CodeMap[model.ResError], http.StatusUnauthorized)
	}
}
func Auth() goweb.HandlerFunc {
	return func(ctx *goweb.Context) {

	}
}
