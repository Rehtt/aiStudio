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
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
	createRecord := &model2.RecordTable{
		Token: info.Token,
		GenID: genId,
	}
	switch req.Type {
	case model.GEN, "":
		createRecord.Prompt = req.Prompt
		if err := sender.GenerateImage(genId, req.Prompt); err != nil {
			logs.Warn("sender.GenerateImage %s", err)
			ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
			return
		}
	case model.VA:
		record, err := repository.GetRecord(req.GenId)
		if err != nil {
			ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
			return
		}
		if record == nil || record.Status != "done" {
			ctx.WriteJSON(model.CodeMap[model.RequestBad], http.StatusBadRequest)
			return
		}
		createRecord.Option = "V" + strconv.Itoa(req.V)
		createRecord.ParentMsgID = record.MsgID
		if err = sender.Variate(genId, int64(req.V), record.MsgID, record.MHash, record.GuildID, record.ChannelID); err != nil {
			logs.Warn("sender.Variate %s", err)
			ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
			return
		}
	case model.UP:
		record, err := repository.GetRecord(req.GenId)
		if err != nil {
			ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
			return
		}
		if record == nil || record.Status != "done" {
			ctx.WriteJSON(model.CodeMap[model.RequestBad], http.StatusBadRequest)
			return
		}
		createRecord.Option = "U" + strconv.Itoa(req.U)
		createRecord.ParentMsgID = record.MsgID
		if err = sender.Upscale(genId, int64(req.V), record.MsgID, record.MHash, record.GuildID, record.ChannelID); err != nil {
			logs.Warn("sender.Upscale %s", err)
			ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
			return
		}
	default:
		ctx.WriteJSON(model.CodeMap[model.RequestBad], http.StatusBadRequest)
		return
	}
	if err := repository.CreateRecord(createRecord); err != nil {
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
		imageUrl = "/api/external/image/" + record.GenID + filepath.Ext(*record.ImageUrl) + "?token=" + record.Token
	}
	ctx.WriteJSON(&model.Response{
		Code: model.OK,
		Data: &model.ProgressResponse{
			Progress: record.Progress,
			Prompt:   record.Prompt,
			Status:   record.Status,
			ImageUrl: imageUrl,
		},
	})
}

func imageUrl(ctx *goweb.Context) {
	id := ctx.GetParam("id")
	if id == "" {
		ctx.WriteJSON(model.CodeMap[model.RequestBad], http.StatusBadRequest)
		return
	}

	record, err := repository.GetRecord(strings.Split(id, ".")[0])
	if err != nil {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
	if record == nil || record.ImageUrl == nil || *record.ImageUrl == "" {
		ctx.WriteJSON(model.CodeMap[model.RequestBad], http.StatusBadRequest)
		return
	}

	req, err := http.Get(*record.ImageUrl)
	if err != nil {
		logs.Warn("http.Get(*record.ImageUrl) error: %s", err)
		return
	}
	defer req.Body.Close()
	ctx.Writer.Header().Set("content-type", req.Header.Get("content-type"))
	ctx.Writer.Header().Set("content-length", req.Header.Get("content-length"))
	ctx.Writer.Header().Set("accept-ranges", req.Header.Get("accept-ranges"))
	ctx.Writer.Header().Set("access-control-allow-origin", "*")
	io.Copy(ctx.Writer, req.Body)
}

func ExternalInfo(ctx *goweb.Context) {
	info, ok := ctx.GetValue("info").(*model.ExternalInfo)
	if !ok {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
	ctx.WriteJSON(model.Response{
		Code: model.OK,
		Data: info,
	})
}
