package user

import (
	"github.com/gin-gonic/gin"
	"hellowiki/common/result"
	"hellowiki/model"
	"hellowiki/service"
	"net/http"
	"strconv"
)

var code int

func Register(c *gin.Context) {
	var userInfo model.RegUser
	_ = c.ShouldBind(&userInfo)
	code = service.CreateUser(&userInfo)
	if code == result.ERROR_USERNAME_USED {
		result.RestFulResult(c, code)
		return
	}
	result.RestFulResult(c, code, userInfo)
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := service.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": result.GetMsg(code),
	})
}

func QueryAllUserInfo(c *gin.Context) {
	var pageSize, pageNum = 10, 1
	pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pageNum, _ = strconv.Atoi(c.Query("pageNum"))
	data := service.GetAllRegUserInfo(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  result.SUCCSE,
		"data":    data,
		"message": result.GetMsg(code),
	})
}

func SetUserName(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var condition model.RegUser
	_ = c.ShouldBind(&condition)
	code = service.SetUser(uint(id), condition)
	if code == result.ERROR_USER_NOT_FOUND {
		result.RestFulResult(c, code)
		return
	}
	result.RestFulResult(c, code, condition)
}
