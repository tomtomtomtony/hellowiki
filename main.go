package main

import (
	"hellowiki/config"
	"hellowiki/routers"
)

func main() {
	//初始化数据库
	config.InitData()
	//初始化路由
	routers.InitRouter().Run(config.Cfg.Server.Port)

}
