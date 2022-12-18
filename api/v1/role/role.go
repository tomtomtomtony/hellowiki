package role

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/result"
	vo2 "hellowiki/api/v1/role/vo"
	"hellowiki/model"
	"hellowiki/service"
	"strconv"
)

// 创建新角色
func RegRole(c *gin.Context) {
	var roleInfo vo2.RoleConditionVO
	_ = c.ShouldBind(&roleInfo)
	code := service.CreateRole(roleInfo)
	result.RestFulResult(c, code)
}

// 从数据库获取角色列表
func QueryAllRoles(c *gin.Context) {
	var pageSize, pageNum = 10, 1
	pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pageNum, _ = strconv.Atoi(c.Query("pageNum"))
	data, total, code := service.GetAllRoles(pageNum, pageSize)
	var res vo2.RoleList
	for i := 0; i < len(data); i++ {
		res.RoleList = append(res.RoleList, do2Vo(data[i]))
	}
	res.Total = total
	result.RestFulResult(c, code, res)
}

func QueryUserRoles(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	res, code := service.GetUserRoles(id)
	result.RestFulResult(c, code, res)
}

func do2Vo(role model.Role) vo2.RoleResult {
	var res vo2.RoleResult
	res.RoleName = role.RoleName
	res.Id = role.ID
	res.CreateAt = role.CreatedAt.UnixMilli()
	res.UpdateAt = role.UpdatedAt.UnixMilli()
	return res
}

func AddPermissionsToUser() {

}

func UpdateUserRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var condition vo2.RoleConditionVO
	_ = c.ShouldBind(&condition)
	code := service.AddRolesForUser(id, condition)
	result.RestFulResult(c, code)

}
