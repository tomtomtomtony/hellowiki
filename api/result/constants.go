package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
*
返回码
*/
const (
	SUCCSE = 200
	ERROR  = 500
	//code=1000....用户模块错误
	ERROR_USERNAME_USED  = 1001
	ERROR_PASSWORD_WRONG = 1002
	ERROR_USER_NOT_FOUND = 1003

	//分类模块错误
	ERROR_CATEGORY_EXIST                   = 2001
	ERROR_CATEGORY_NOT_FOUND               = 2002
	ERROR_PARENT_CATEGORY_NOT_FOUND        = 2003
	ERROR_ARTICLE_DATABASE_Index_NOT_FOUND = 2004

	//文章模块

)

var codeMsg = map[int]string{
	SUCCSE:                                 "OK",
	ERROR:                                  "FAIL",
	ERROR_USERNAME_USED:                    "用户已存在",
	ERROR_USER_NOT_FOUND:                   "用户不存在",
	ERROR_CATEGORY_EXIST:                   "该分类已存在",
	ERROR_CATEGORY_NOT_FOUND:               "该分类不存在",
	ERROR_PARENT_CATEGORY_NOT_FOUND:        "指定的上级菜单不存在",
	ERROR_ARTICLE_DATABASE_Index_NOT_FOUND: "数据库错误:分类存储表未找到",
}

func GetMsg(code int) string {
	return codeMsg[code]
}

func RestFulResult(c *gin.Context, statuesCode int, data ...interface{}) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "" {
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
	}
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	if data != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  statuesCode,
			"data":    data,
			"message": GetMsg(statuesCode),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  statuesCode,
		"message": GetMsg(statuesCode),
	})

}
