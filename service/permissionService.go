package service

import (
	"hellowiki/api/result"
	"hellowiki/api/v1/permission/vo"
	"hellowiki/model"
	"log"
)

func AddPermission(condition vo.PermissionVO) (code int) {
	permissionInfo := permissionVO2Do(condition)
	if model.HasPermissionByName(permissionInfo.PermissionName) != result.SUCCSE {
		log.Printf("权限已存在")
		return result.ERROR_PERMISSION_EXIST
	}
	return model.InsertPermissionInfo(permissionInfo)
}

func GetAllPermission(pageNum int, pageSize int) (res []model.Permission, total int64, code int) {
	res, total = model.GetAllPermissionType(pageNum, pageSize)
	return res, total, result.SUCCSE
}

func permissionVO2Do(vo vo.PermissionVO) model.Permission {
	var res model.Permission
	res.PermissionName = vo.PermissionName
	return res
}
