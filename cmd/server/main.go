package main

import (
	"aiStudio/internal/conf"
	"aiStudio/internal/midj"
	"aiStudio/internal/mysql"
	"aiStudio/internal/redis"
	"aiStudio/internal/service"
	"flag"
	"github.com/Rehtt/Kit/i18n"
	"github.com/Rehtt/Kit/log/logs"
	goweb "github.com/Rehtt/Kit/web"
	"net/http"
	"os"
)

var (
	config                 = flag.String("conf", "configs", i18n.GetText("配置文件位置"))
	generateConfigTemplate = flag.Bool("gen-conf", false, "生成配置文件模板")
)

func main() {
	flag.Parse()

	if *generateConfigTemplate {
		if err := conf.GenConfig(*config); err != nil {
			logs.Fatal(i18n.GetText("模板生成错误：%s"), err)
		}
		os.Exit(0)
	}

	conf.Init(*config)
	if err := redis.New(&conf.GetServer().Redis); err != nil {
		logs.Fatal(i18n.GetText("Redis 初始化失败：%s"), err)
	}
	logs.Info(i18n.GetText("Redis 初始化成功"))

	if err := mysql.Init(&conf.GetServer().Mysql); err != nil {
		logs.Fatal(i18n.GetText("Mysql 初始化失败：%s"), err)
	}
	logs.Info(i18n.GetText("Mysql 初始化成功"))

	if err := midj.Init(&conf.GetServer().Midj); err != nil {
		logs.Fatal(i18n.GetText("Midj 初始化失败：%s"), err)
	}
	logs.Info(i18n.GetText("Midj 初始化成功"))

	web := goweb.New()
	service.Route(web)
	http.ListenAndServe(conf.GetServer().Listen, web)
}
