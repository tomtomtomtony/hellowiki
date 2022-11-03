package routers

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/v1/category"
	"hellowiki/api/v1/user"
	"hellowiki/config"
)

func InitRouter() *gin.Engine {
	gin.SetMode(config.Cfg.Server.AppMode)

	r := gin.Default()
	routerUserV1 := r.Group("api/v1/user")
	{
		//用户模块
		routerUserV1.POST("/register", user.Register)
		routerUserV1.GET("/all", user.QueryAllUserInfo)
		routerUserV1.DELETE("/del/:id", user.DeleteUser)
		routerUserV1.PUT("/edt/:id", user.SetUserName)

	}

	routerCategoryV1 := r.Group("api/v1/category")
	{
		//分类模块
		routerCategoryV1.POST("/create", category.CreateCategory)
		routerCategoryV1.DELETE("/del/:id", category.DeleteCategory)
		routerCategoryV1.GET("/all", category.QueryAllCategory)
		routerCategoryV1.PUT("/rename/:id", category.ReNameCategory)
	}
	return r
}
