package category

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/result"
	"hellowiki/api/v1/category/vo"
	"hellowiki/service"
)

var code int

// 创建分类
func CreateCategory(c *gin.Context) {
	var condition vo.CategoryVO
	_ = c.ShouldBind(&condition)

	code = service.CreateCategory(condition)
	result.RestFulResult(c, code)
}

/*
*
删除指定分类，需要提供分类的id,parentId,parentName,engName
*/
func DeleteCategory(c *gin.Context) {
	var condition vo.CategoryVO
	_ = c.ShouldBind(&condition)
	code := service.DeleteCategory(condition)
	result.RestFulResult(c, code)
}

// 重建对应分类的数据库表。用于content有效，但数据库对应分类文章表损坏的情况
func repairTable(c *gin.Context) {

}

func GetAllTopCategory(c *gin.Context) {
	data := service.GetTopCategory()
	result.RestFulResult(c, result.SUCCSE, data)
}

func GetNextLevelMenu(c *gin.Context) {
	currentPath := c.Query("path")
	data := service.GetDirectChildren(currentPath)
	result.RestFulResult(c, result.SUCCSE, data)
}
func GetAllArticle(c *gin.Context) {
	var condition vo.CategoryVO
	_ = c.ShouldBind(&condition)
	detail, total := service.GetArticleList(condition)
	result.RestFulResult(c, result.SUCCSE, detail, total)
}
