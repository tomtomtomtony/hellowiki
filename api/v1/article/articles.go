package article

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/v1/article/vo"
	"hellowiki/common/result"
	"hellowiki/service"
)

// 创建文章
func CreateArticle(c *gin.Context) {
	var article vo.ConditionVO
	_ = c.ShouldBind(&article)
	code := service.CreateArticle(article)
	result.RestFulResult(c, code)
}

func GetAllInCategory(c *gin.Context) {
	var condition vo.ConditionVO
	_ = c.BindJSON(&condition)
	data, code := service.QueryInCategory(condition)
	if code != result.SUCCSE {
		result.RestFulResult(c, code)
	}
	result.RestFulResult(c, result.SUCCSE, data)
}
