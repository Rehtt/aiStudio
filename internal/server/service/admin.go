package service

import (
	"aiStudio/internal/redis"
	"aiStudio/internal/server/service/auth"
	"aiStudio/internal/server/service/model"
	"github.com/Rehtt/Kit/log/logs"
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
)

func SetUser(ctx *goweb.Context) {
	var user model.AddUser
	err := ctx.ReadJSON(&user)
	if err != nil {
		logs.Warn("AddUser ctx.ReadJSON(&user) err: %s", err)
		ctx.WriteJSON(model.CodeMap[model.RequestBad], http.StatusBadRequest)
		return
	}
	err = redis.DB.Set(ctx, auth.ExternalAuthKey+user.Token, user.Number, user.ExpDuration).Err()
	if err != nil {
		logs.Warn("AddUser redis.DB.Set err: %s", err)
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
	ctx.WriteJSON(model.Response{
		Code: model.OK,
	})
}
