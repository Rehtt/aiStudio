package main

import (
	"aiStudio/internal/conf"
	_const "aiStudio/internal/const"
	"aiStudio/internal/midj/listener"
	"aiStudio/internal/redis"
	"context"
	"flag"
	"github.com/Rehtt/Kit/i18n"
	"github.com/Rehtt/Kit/log/logs"
	jsoniter "github.com/json-iterator/go"
)

var (
	config = flag.String("conf", "configs", i18n.GetText("配置文件位置"))
)

func main() {
	flag.Parse()
	conf.Init(*config)
	if err := redis.New(&conf.GetServer().Redis); err != nil {
		logs.Fatal(i18n.GetText("Redis 初始化失败：%s"), err)
	}
	logs.Info(i18n.GetText("Redis 初始化成功"))

	// midj
	if err := listener.Init(conf.GetServer().Midj); err != nil {
		logs.Fatal(i18n.GetText("Midj 初始化失败：%s"), err)
	}
	logs.Info(i18n.GetText("Midj 初始化成功"))
	// 收到消息推送到队列
	listener.Request(func(r listener.ReqCb) {
		d, _ := jsoniter.Marshal(r)
		err := redis.DB.LPush(context.Background(), _const.DISCORD_RESULTS_QUEUE, d).Err()
		if err != nil {
			panic(err)
		}
		logs.Info("push %s", r.Type)
	})

	<-make(chan struct{})
}
