package main

import (
	"hellowiki/config"
	"hellowiki/model"
	"hellowiki/routers"
)

func main() {
	//初始化数据库
	model.InitData()
	//初始化路由
	routers.InitRouter().Run(config.Cfg.Server.Port)

}
