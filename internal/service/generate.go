package service

import (
	"aiStudio/internal/midj"
	"aiStudio/internal/redis"
	"aiStudio/internal/service/model"
	"aiStudio/pkg"
	"github.com/Rehtt/Kit/log/logs"
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
)

func generate(ctx *goweb.Context) {
	var req model.GenerateImageRequest
	if err := ctx.ReadJSON(&req); err != nil || req.Prompt == "" {
		ctx.WriteJSON(model.CodeMap[model.RequestBad], http.StatusBadRequest)
		return
	}
	info, ok := ctx.GetValue("info").(*model.ExternalInfo)
	if !ok {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
	defer func() {
		redis.UnLock(ctx, info.LockKey)
	}()
	if err := redis.DB.Set(ctx, info.Key, info.Number-1, info.ExpirationDuration).Err(); err != nil {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
	if err := midj.GenerateImage(req.Prompt); err != nil {
		logs.Warn("%s", err)
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}

	ctx.WriteJSON(&model.Response{
		Code: model.OK,
		Data: &model.GenerateImageResponse{Id: pkg.GenId(info)},
	})
}
