package service

import (
	"aiStudio/internal/service/auth"
	goweb "github.com/Rehtt/Kit/web"
)

func Route(g *goweb.GOweb) {
	var (
		api      = g.Grep("/api")
		external = api.Grep("/external")
	)
	// 外部
	{
		external.Middleware(auth.ExternalAuth())
		external.FootMiddleware(auth.EUnlock)

		external.POST("/generate", generate)
		external.GET("/progress", progress)
	}

	// 后台
	{
		api.Middleware(auth.Auth())

	}
}
