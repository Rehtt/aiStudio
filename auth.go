package main

import (
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
)

func Auth() goweb.HandlerFunc {
	return func(ctx *goweb.Context) {
		token := ctx.Request.URL.Query().Get("token")
		if token == "" {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			ctx.Stop()
			return
		}

	}
}
