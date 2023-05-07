package main

import (
	goweb "github.com/Rehtt/Kit/web"
	"github.com/redis/go-redis/v9"
)

var r *redis.Client

func Redis(ctx *goweb.GOweb, addr, password string, db int) error {
	r = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	_, err := r.Ping(ctx).Result()
	ctx.SetValue("redis", r)
	return err
}
