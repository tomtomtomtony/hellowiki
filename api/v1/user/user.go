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
		code = result.ERROR_USERNAME_USED
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    userInfo,
		"message": result.GetErrMsg(code),
	})
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := service.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": result.GetErrMsg(code),
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
		"message": result.GetErrMsg(code),
	})
}

func SetUserName(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var condition model.RegUser
	_ = c.ShouldBind(&condition)
	code = service.SetUserName(uint(id), condition)
	if code == result.ERROR_USER_NOT_FOUND {
		code = result.ERROR_USER_NOT_FOUND
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": result.GetErrMsg(code),
	})
}
