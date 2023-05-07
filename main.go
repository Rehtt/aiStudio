package main

import (
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
)

func main() {
	web := goweb.New()
	Route(web)
	if err := Redis(web, "", "", 0); err != nil {

	}
	http.ListenAndServe(":8080", web)
}
