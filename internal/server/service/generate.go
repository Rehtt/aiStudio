package service

import (
	"aiStudio/internal/midj/sender"
	"aiStudio/internal/redis"
	"aiStudio/internal/repository"
	model2 "aiStudio/internal/repository/model"
	model "aiStudio/internal/server/service/model"
	"aiStudio/pkg"
	"github.com/Rehtt/Kit/log/logs"
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

	genId := pkg.GenId(info, req)
	if err := repository.CreateRecord(&model2.RecordTable{
		Token:    info.Token,
		Prompt:   req.Prompt,
		GenID:    genId,
		Progress: 0,
		//ParentMsgID: nil,
		//Option:      nil,
	}); err != nil {
		logs.Warn("%s", err)
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}

	if err := sender.GenerateImage(genId, req.Prompt); err != nil {
		logs.Warn("%s", err)
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}

	ctx.WriteJSON(&model.Response{
		Code: model.OK,
		Data: &model.GenerateImageResponse{Id: genId},
	})
}

func progress(ctx *goweb.Context) {
	id := ctx.GetParam("id")
	record, err := repository.GetRecord(id)
	if err != nil {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
	if record == nil {
		ctx.WriteJSON(model.CodeMap[model.RequestBad], http.StatusBadRequest)
		return
	}
	var imageUrl string
	if record.ImageUrl != nil {
		imageUrl = *record.ImageUrl
	}
	ctx.WriteJSON(&model.Response{
		Code: model.OK,
		Data: &model.ProgressResponse{
			Progress: record.Progress,
			ImageUrl: imageUrl,
		},
	})
}
