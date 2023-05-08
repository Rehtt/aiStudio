package service

import (
	"aiStudio/internal/redis"
	"aiStudio/internal/service/model"
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
)

func generate(ctx *goweb.Context) {
	info, ok := ctx.GetValue("info").(*model.ExternalInfo)
	if !ok {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
	if err := redis.DB.Set(ctx, info.Key, info.Number-1, info.ExpirationDuration).Err(); err != nil {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
}
