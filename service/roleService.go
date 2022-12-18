package service

import (
	"hellowiki/api/result"
	vo2 "hellowiki/api/v1/role/vo"
	utils2 "hellowiki/common/utils"
	"hellowiki/model"
	"log"
	"strconv"
)

func GetAllRoles(pageNum int, pageSize int) (res []model.Role, total int64, code int) {
	res, total = model.GetAllRoles(pageNum, pageSize)
	return res, total, result.SUCCSE
}
func GetUserRoles(userId int) (res []string, code int) {
	if model.HasUserById(uint(userId)) == result.SUCCSE {
		return res, result.ERROR_USER_NOT_FOUND
	}
	return model.GetRolesForUserById(userId)
}

func AddRolesForUser(userId int, condition vo2.RoleConditionVO) (code int) {
	if model.HasUserById(uint(userId)) == result.SUCCSE {
		return result.ERROR_USER_NOT_FOUND
	}

	enforcer := utils2.GetEnforcer()
	userIdStr := strconv.Itoa(userId)
	//清空当前用户的角色
	curr, code := GetUserRoles(userId)
	if code != result.SUCCSE {
		log.Printf("未能正确获取当前用户角色")
	}
	for _, role := range curr {
		_, err := enforcer.DeleteRoleForUser(userIdStr, role)
		if err != nil {
			log.Printf("未能正确删除角色{%v}:{%v}", role, err)
			return 0
		}
	}
	//添加目标角色
	for _, role := range condition.Roles {
		enforcer.AddRoleForUser(userIdStr, role)
	}
	err := enforcer.SavePolicy()
	if err != nil {
		log.Printf("保持策略时发生错误:{%v}", err)
		return result.ERROR
	}
	return result.SUCCSE
}

func CreateRole(condition vo2.RoleConditionVO) (code int) {
	var roleInfo = roleVo2Do(condition)
	codeHas := model.HasRoleByName(roleInfo.RoleName)
	if codeHas != result.SUCCSE {
		return codeHas
	}
	codeInsert := model.InsertRoleInfo(roleInfo)
	if codeInsert != result.SUCCSE {
		return codeInsert
	}
	return codeInsert
}

func roleVo2Do(vo vo2.RoleConditionVO) model.Role {
	var do model.Role
	do.RoleName = vo.RoleName
	return do
}
