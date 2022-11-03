package article

import (
	"github.com/gin-gonic/gin"
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
