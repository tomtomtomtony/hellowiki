package service

import (
	"hellowiki/api/result"
	vo2 "hellowiki/api/v1/role/vo"
	"hellowiki/model"
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
