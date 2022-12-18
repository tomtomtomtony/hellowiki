package permission

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/result"
	"hellowiki/api/v1/permission/vo"
	"hellowiki/model"
	"hellowiki/service"
	"strconv"
)

// 创建新权限
func CreatePermission(c *gin.Context) {
	var permissionInfo vo.PermissionVO
	_ = c.ShouldBind(&permissionInfo)
	code := service.AddPermission(permissionInfo)
	result.RestFulResult(c, code)
}

//获取当前数据库全部权限类型

// 从数据库获取角色列表
func QueryAllPermission(c *gin.Context) {
	var pageSize, pageNum = 10, 1
	pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pageNum, _ = strconv.Atoi(c.Query("pageNum"))
	data, total, code := service.GetAllPermission(pageNum, pageSize)
	var res vo.PermissionList
	for i := 0; i < len(data); i++ {
		res.PermissionList = append(res.PermissionList, do2Vo(data[i]))
	}
	res.Total = total
	result.RestFulResult(c, code, res)
}
func do2Vo(do model.Permission) vo.PermissionResult {
	var res vo.PermissionResult
	res.PermissionName = do.PermissionName
	res.Id = do.ID
	res.CreateAt = do.CreatedAt.UnixMilli()
	res.UpdateAt = do.UpdatedAt.UnixMilli()
	return res
}
