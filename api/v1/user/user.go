package user

import (
	"github.com/gin-gonic/gin"
	"hellowiki/common/result"
	"hellowiki/model"
	"hellowiki/service"
	"net/http"
)

var code int

func Register(c *gin.Context) {
	var userInfo model.RegUser
	_ = c.ShouldBind(&userInfo)
	code = service.CreateUser(&userInfo)
	if code == result.ERROR_USERNAME_USED {
		code = result.ERROR_USERNAME_USED
	}
	c.JSONP(http.StatusOK, gin.H{
		"status":  code,
		"data":    userInfo,
		"message": result.GetErrMsg(code),
	})
}
