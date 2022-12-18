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

func Login(c *gin.Context) {
	var userInfo vo.LoginUserVO
	_ = c.ShouldBind(&userInfo)
	userName, token, code := service.UserLogin(userInfo)
	var res vo.LoginUserVO
	res.Token = token
	res.UserName = userName
	res.Code = code
	result.RestFulResult(c, code, res)
}

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
	data, total := service.GetAllRegUserInfo(pageSize, pageNum)
	var res vo.UserList
	for i := 0; i < len(data); i++ {
		res.UserList = append(res.UserList, do2ResultVo(data[i]))
	}
	res.Total = total
	result.RestFulResult(c, result.SUCCSE, res)
}

func SetUserName(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var condition vo.RegUserVO
	_ = c.ShouldBind(&condition)
	code = service.SetUserName(uint(id), condition)
	if code == result.ERROR_USER_NOT_FOUND {
		result.RestFulResult(c, code)
		log.Printf("用户不存在")
	}
	result.RestFulResult(c, code)
}

func SetUserRoles(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var condition vo.RegUserVO
	_ = c.ShouldBind(&condition)
	code = service.SetRoles(uint(id), condition)
	if code == result.ERROR_USER_NOT_FOUND {
		result.RestFulResult(c, code)
		log.Printf("用户不存在")
	}
	result.RestFulResult(c, code)
}

func do2ResultVo(regUser model.RegUser) vo.UserResult {
	var res vo.UserResult
	res.UserName = regUser.UserName
	res.Id = regUser.ID
	res.CreateAt = regUser.CreatedAt.UnixMilli()
	res.UpdateAt = regUser.UpdatedAt.UnixMilli()
	return res
}
