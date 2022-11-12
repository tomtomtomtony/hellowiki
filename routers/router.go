package routers

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/v1/article"
	"hellowiki/api/v1/category"
	"hellowiki/api/v1/user"
)

func InitRouter() *gin.Engine {

	r := gin.Default()
	//用户模块
	routerUserV1 := r.Group("api/v1/user")
	{
		routerUserV1.POST("/register", user.Register)
		routerUserV1.GET("/all", user.QueryAllUserInfo)
		routerUserV1.DELETE("/del/:id", user.DeleteUser)
		routerUserV1.PUT("/edt/:id", user.SetUserName)

	}
	//分类模块
	routerCategoryV1 := r.Group("api/v1/category")
	{
		routerCategoryV1.POST("/create", category.CreateCategory)
		routerCategoryV1.DELETE("/del/:id", category.DeleteCategory)
		routerCategoryV1.GET("/all", category.QueryAllCategory)
		routerCategoryV1.PUT("/rename/:id", category.ReNameCategory)
	}

	routerArticleV1 := r.Group("api/v1/article")
	{
		routerArticleV1.POST("/create", article.CreateArticle)
		routerArticleV1.POST("/allInCategory", article.GetAllTitleCurrentCategory)
	}

	return r
}
