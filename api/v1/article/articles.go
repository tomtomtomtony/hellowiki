package article

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/result"
	"hellowiki/api/v1/article/vo"
	"hellowiki/service"
)

// 创建文章
func CreateArticle(c *gin.Context) {
	var article vo.ConditionVO
	_ = c.ShouldBind(&article)
	code := service.CreateArticle(article)
	result.RestFulResult(c, code)
}

// 获取一个分类下所有文章
func GetAllTitleCurrentCategory(c *gin.Context) {
	var condition vo.ConditionVO
	_ = c.BindJSON(&condition)
	data, code := service.QueryInCategory(condition)
	if code != result.SUCCSE {
		result.RestFulResult(c, code)
	}
	result.RestFulResult(c, result.SUCCSE, data)
}

// 获取指定的文章
func GetArticle(c *gin.Context) {
	var article vo.ConditionVO
	_ = c.ShouldBind(&article)
	res, code := service.GetArticle(article)
	if code != result.SUCCSE {
		result.RestFulResult(c, code)
	}
	result.RestFulResult(c, code, res)
}

func DelArticle(c *gin.Context) {
	var article vo.ConditionVO
	_ = c.ShouldBind(&article)
	code := service.DeleteArticle(article)
	result.RestFulResult(c, code)
}
