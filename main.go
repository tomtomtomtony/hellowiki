package main

import (
	"github.com/gin-gonic/gin"
	"hellowiki/config"
	"hellowiki/routers"
)

func main() {
	gin.SetMode(config.Cfg.Server.AppMode)
	//初始化路由
	routers.InitRouter().Run(config.Cfg.Server.Port)

}
