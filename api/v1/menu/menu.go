package menu

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/result"
	"hellowiki/service"
	"strconv"
)

func GetAllMenuCurrentCategory(c *gin.Context) {
	categoryId, _ := strconv.Atoi(c.Query("categoryId"))

	data := service.GetDirectChildren(uint(categoryId))
	result.RestFulResult(c, result.SUCCSE, data)
}

func GetAllTopCategory(c *gin.Context) {
	data := service.GetAllTopCategory()
	result.RestFulResult(c, result.SUCCSE, data)
}
