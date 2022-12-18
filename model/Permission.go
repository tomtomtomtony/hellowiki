package model

import (
	"gorm.io/gorm"
	"hellowiki/api/result"
	utils2 "hellowiki/common/utils"
	"log"
)

// 设计读写等权限
type Permission struct {
	gorm.Model
	PermissionName string `gorm:"type:varchar(32);not null" json:"permissionName"`
}

func GetAllPermissionType(pageNum int, pageSize int) (permissions []Permission, total int64) {
	dbBase := utils2.OpenDB()
	err := dbBase.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&permissions).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return permissions, 0
	}
	return permissions, total
}

// 插入用户数据
func InsertPermissionInfo(data Permission) (code int) {
	dbBase := utils2.OpenDB()
	err := dbBase.Create(&data).Error
	if err != nil {
		log.Printf("插入权限错误:{%v}", err)
		return result.ERROR
	}
	return result.SUCCSE
}

func HasPermissionByName(permissionName string) (code int) {
	var permission Permission
	dbBase := utils2.OpenDB()
	dbBase.Take(&permission, "permission_name=?", permissionName)
	if permission.ID > 0 {
		//权限已存在
		return result.ERROR_USERNAME_USED
	}
	//权限不存在
	return result.SUCCSE
}
