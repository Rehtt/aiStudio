package auth

import (
	"aiStudio/internal/redis"
	model2 "aiStudio/internal/server/service/model"
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
)

// EUnlock 解锁
func EUnlock(ctx *goweb.Context) {
	info, ok := ctx.GetValue("info").(*model2.ExternalInfo)
	if !ok {
		ctx.WriteJSON(model2.CodeMap[model2.ServerBad], http.StatusBadGateway)
		return
	}
	if err := redis.UnLock(ctx, info.LockKey); err != nil {
		ctx.WriteJSON(model2.CodeMap[model2.ServerBad], http.StatusBadGateway)
		return
	}
}
