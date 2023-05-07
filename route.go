package main

import goweb "github.com/Rehtt/Kit/web"

func Route(g *goweb.GOweb) {
	var (
		api      = g.Grep("/api")
		external = api.Grep("/external")
	)
	external.Middleware(Auth())
}
