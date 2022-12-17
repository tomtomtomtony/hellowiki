package routers

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/v1/article"
	"hellowiki/api/v1/category"
	"hellowiki/api/v1/role"
	"hellowiki/api/v1/user"
	"hellowiki/middleware"
)

func InitRouter() *gin.Engine {

	r := gin.Default()
	r.Use(middleware.CrossHandler())
	//用户模块
	routerUserV1 := r.Group("api/v1/user")
	{
		routerUserV1.POST("/register", user.Register)
		routerUserV1.POST("/login", user.Login)
		routerUserV1.GET("/all", user.QueryAllUserInfo)
		routerUserV1.DELETE("/del/:id", user.DeleteUser)
		routerUserV1.PUT("/editName/:id", user.SetUserName)
		routerUserV1.PUT("/editRole/:id", user.SetUserRoles)
	}
	//角色模块
	routerRoleV1 := r.Group("api/v1/role")
	{
		routerRoleV1.GET("/all", role.QueryAllRoles)
		routerRoleV1.GET("/getRoles", role.QueryUserRoles)
		routerRoleV1.POST("/create", role.RegRole)
	}

	//分类模块
	routerCategoryV1 := r.Group("api/v1/category")
	routerCategoryV1.Use(middleware.JWTAuth())
	{
		routerCategoryV1.POST("/create", category.CreateCategory)
		routerCategoryV1.POST("/del", category.DeleteCategory)
		routerCategoryV1.GET("/getTop", category.GetAllTopCategory)
		routerCategoryV1.GET("/currentall", category.GetNextLevelMenu)
	}

	routerArticleV1 := r.Group("api/v1/article")
	routerArticleV1.Use(middleware.JWTAuth())
	{
		routerArticleV1.POST("/create", article.CreateArticle)
		routerArticleV1.POST("/allInCategory", article.GetAllTitleCurrentCategory)
		routerArticleV1.POST("/getArticle", article.GetArticle)
		routerArticleV1.POST("/del", article.DelArticle)
		routerArticleV1.POST("/update", article.UpdateArticle)
	}

	return r
}
