package service

import (
	"aiStudio/internal/redis"
	"aiStudio/internal/service/model"
	"aiStudio/pkg"
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
)

func generate(ctx *goweb.Context) {
	info, ok := ctx.GetValue("info").(*model.ExternalInfo)
	if !ok {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}

	var req model.GenerateImageRequest
	if err := ctx.ReadJSON(&req); err != nil || req.Prompt == "" {
		ctx.WriteJSON(model.CodeMap[model.RequestBad], http.StatusBadRequest)
		return
	}
	if info.Number <= 0 {
		ctx.WriteJSON(model.CodeMap[model.ResError], http.StatusUnauthorized)
		return
	}

	// 次数-1
	if err := redis.DB.Decr(ctx, info.RedisKey).Err(); err != nil {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}

	//if err := repository.CreateRecord(info.Token, req.Prompt, pkg.GenId(info)); err != nil {
	//	logs.Warn("%s", err)
	//	ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
	//	return
	//}

	//if err := midj.GenerateImage(req.Prompt); err != nil {
	//	logs.Warn("%s", err)
	//	ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
	//	return
	//}

	ctx.WriteJSON(&model.Response{
		Code: model.OK,
		Data: &model.GenerateImageResponse{Id: pkg.GenId(info)},
	})
}

func progress(ctx *goweb.Context) {

}
