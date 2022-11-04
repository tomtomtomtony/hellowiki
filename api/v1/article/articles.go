package article

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/v1/article/vo"
	"hellowiki/common/result"
	"hellowiki/model"
	"hellowiki/service"
)

// 创建文章
func CreateArticle(c *gin.Context) {
	var article model.Article
	_ = c.ShouldBind(&article)
	code := service.CreateArticle(article)
	result.RestFulResult(c, code)
}

func GetAllInCategory(c *gin.Context) {
	var condition vo.ConditionVO
	_ = c.BindJSON(&condition)
	data := service.QueryInCategory(condition)
	if data == nil {
		result.RestFulResult(c, result.ERROR)
		return
	}
	result.RestFulResult(c, result.SUCCSE, data)
}
