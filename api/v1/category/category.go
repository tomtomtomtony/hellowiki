package category

import (
	"github.com/gin-gonic/gin"
	"hellowiki/common/result"
	"hellowiki/model"
	"hellowiki/service"
	"net/http"
	"strconv"
)

var code int

// 创建分类
func CreateCategory(c *gin.Context) {
	var category model.Category
	_ = c.ShouldBind(&category)
	if category.ParentId == 0 {
		code = service.CreateRootCategory(category.Name)
	} else {
		code = service.CreateNonRootCategory(category)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    category,
		"message": result.GetErrMsg(code),
	})
}

func DeleteCategory(c *gin.Context) {
	var category model.Category
	_ = c.ShouldBind(&category)
	code := service.DeleteCategory(category)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": result.GetErrMsg(code),
	})
}

func QueryAllCategory(c *gin.Context) {
	var pageSize, pageNum = 10, 1
	pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pageNum, _ = strconv.Atoi(c.Query("pageNum"))
	data := service.GetAllCategory(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  result.SUCCSE,
		"data":    data,
		"message": result.GetErrMsg(code),
	})
}

func ReNameCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	newName := c.Param("name")
	var condition model.Category
	condition.Name = newName
	code = service.SetCategory(uint(id), condition)
	if code == result.ERROR_USER_NOT_FOUND {
		code = result.ERROR_USER_NOT_FOUND
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": result.GetErrMsg(code),
	})
}
