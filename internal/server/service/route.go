package service

import (
	auth2 "aiStudio/internal/server/service/auth"
	goweb "github.com/Rehtt/Kit/web"
)

func Route(g *goweb.GOweb) {
	var (
		api      = g.Grep("/api")
		external = api.Grep("/external")
	)
	// 外部
	{
		external.Middleware(auth2.ExternalAuth())
		external.FootMiddleware(auth2.EUnlock)

		external.POST("/generate", generate)
		external.GET("/progress/#id", progress)
	}

	// 后台
	{
		api.Middleware(auth2.Auth())

	}
}
