package user

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/result"
	"hellowiki/api/v1/user/vo"
	"hellowiki/model"
	"hellowiki/service"
	"log"
	"strconv"
)

var code int

func Register(c *gin.Context) {
	var userInfo vo.RegUserVO
	_ = c.ShouldBind(&userInfo)
	code = service.CreateUser(userInfo)
	if code == result.ERROR_USERNAME_USED {
		result.RestFulResult(c, code)
		log.Printf("用户名已经被使用")
	}
	result.RestFulResult(c, code, userInfo)
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := service.DeleteUser(id)
	result.RestFulResult(c, code)
}

func QueryAllUserInfo(c *gin.Context) {
	var pageSize, pageNum = 10, 1
	pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pageNum, _ = strconv.Atoi(c.Query("pageNum"))
	data := service.GetAllRegUserInfo(pageSize, pageNum)
	result.RestFulResult(c, result.SUCCSE, data)
}

func SetUserName(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var condition model.RegUser
	_ = c.ShouldBind(&condition)
	code = service.SetUser(uint(id), condition)
	if code == result.ERROR_USER_NOT_FOUND {
		result.RestFulResult(c, code)
		log.Fatalf("用户不存在")
	}
	result.RestFulResult(c, code, condition)
}
