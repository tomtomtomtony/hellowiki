package routers

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/v1/user"
	"hellowiki/config"
)

func InitRouter() *gin.Engine {
	gin.SetMode(config.Cfg.Server.AppMode)

	r := gin.Default()
	routerV1 := r.Group("api/v1")
	{
		routerV1.POST("user/register", user.Register)
	}
	return r
}
