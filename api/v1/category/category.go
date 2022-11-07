package category

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/result"
	"hellowiki/api/v1/category/vo"
	"hellowiki/model"
	"hellowiki/service"
	"strconv"
)

var code int

// 创建分类
func CreateCategory(c *gin.Context) {
	var condition vo.ConditionVO
	_ = c.ShouldBind(&condition)
	code = service.CreateCategory(condition)
	result.RestFulResult(c, code)
}

/*
*
删除指定分类，需要提供分类的id,parentId,parentName,engName
*/
func DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var category model.Category
	category.ID = uint(id)
	_ = c.ShouldBind(&category)
	code := service.DeleteCategory(category)
	result.RestFulResult(c, code)

}

func QueryAllCategory(c *gin.Context) {
	var pageSize, pageNum = 10, 1
	pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pageNum, _ = strconv.Atoi(c.Query("pageNum"))
	data := service.GetAllCategory(pageSize, pageNum)
	if data == nil {
		result.RestFulResult(c, result.ERROR)
		return
	}
	result.RestFulResult(c, result.SUCCSE, data)
}

func ReNameCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var condition model.Category
	_ = c.ShouldBind(&condition)
	code = service.SetCategory(uint(id), condition)
	result.RestFulResult(c, code)
}

// 重建对应分类的数据库表。用于分类存在，但数据库对应分类文章表损坏的情况
func repairTable(c *gin.Context) {

}
