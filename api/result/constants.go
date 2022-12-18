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
	ERROR_CATEGORY_NOT_EMPTY               = 2005
	//文章模块
	ERROR_ARTICLE_NOT_FOUND = 3001
	//jwt错误
	ERROR_TOKEN_EXPIRED            = 4001
	ERROR_TOKEN_EXPIRED_MaxRefresh = 4002
	ERROR_TOKEN_MALFORMED          = 4003
	ERROR_TOKEN_INVALID            = 4004
	ERROR_HEADER_EMPTY             = 4005
	ERROR_HEADER_MALFORMED         = 4006
	ERROR_TOKEN_NOT_FOUND          = 4007

	//角色模块
	ERROR_ROLE_EXIST = 5001

	//权限模块
	ERROR_PERMISSION_EXIST = 6001
)

var codeMsg = map[int]string{
	SUCCSE:                                 "操作成功",
	ERROR:                                  "FAIL",
	ERROR_USERNAME_USED:                    "用户已存在",
	ERROR_USER_NOT_FOUND:                   "用户不存在",
	ERROR_CATEGORY_EXIST:                   "该分类已存在",
	ERROR_CATEGORY_NOT_FOUND:               "该分类不存在",
	ERROR_PARENT_CATEGORY_NOT_FOUND:        "指定的上级菜单不存在",
	ERROR_ARTICLE_DATABASE_Index_NOT_FOUND: "数据库错误:分类存储表未找到",
	ERROR_TOKEN_EXPIRED:                    "令牌已过期",
	ERROR_TOKEN_EXPIRED_MaxRefresh:         "令牌已过最大刷新时间",
	ERROR_TOKEN_MALFORMED:                  "请求令牌格式有误",
	ERROR_TOKEN_INVALID:                    "请求令牌无效",
	ERROR_HEADER_EMPTY:                     "需要认证才能访问",
	ERROR_HEADER_MALFORMED:                 "请求头中 Authorization 格式有误",
	ERROR_TOKEN_NOT_FOUND:                  "请求未携带令牌，无权限访问",
	ERROR_ARTICLE_NOT_FOUND:                "指定文章未找到",
	ERROR_CATEGORY_NOT_EMPTY:               "分类下不为空",
	ERROR_ROLE_EXIST:                       "角色已存在",
	ERROR_PERMISSION_EXIST:                 "权限已存在",
}

func GetMsg(code int) string {
	return codeMsg[code]
}

func RestFulResult(c *gin.Context, statuesCode int, data ...interface{}) {
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
