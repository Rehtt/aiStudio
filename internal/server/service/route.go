package service

import (
	auth2 "aiStudio/internal/server/service/auth"
	goweb "github.com/Rehtt/Kit/web"
)

func Route(g *goweb.GOweb) {
	var (
		api      = g.Grep("/api")
		external = api.Grep("/external")
		admin    = api.Grep("/admin")
	)
	// 外部
	{
		external.Middleware(auth2.ExternalAuth())

		external.GET("/progress/#id", progress)
		external.GET("/image/#id", imageUrl)

		external.FootMiddleware(auth2.EUnlock)
		external.POST("/generate", generate)
		external.GET("/info", ExternalInfo)
	}

	// 后台
	{
		admin.Middleware(auth2.Auth())
		admin.POST("/set/user", SetUser)
	}
}
