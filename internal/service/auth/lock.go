package auth

import (
	"aiStudio/internal/redis"
	"aiStudio/internal/service/model"
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
)

// EUnlock 解锁
func EUnlock(ctx *goweb.Context) {
	info, ok := ctx.GetValue("info").(*model.ExternalInfo)
	if !ok {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
	if err := redis.UnLock(ctx, info.LockKey); err != nil {
		ctx.WriteJSON(model.CodeMap[model.ServerBad], http.StatusBadGateway)
		return
	}
}
