package model

import (
	"gorm.io/gorm"
	"hellowiki/api/result"
	utils2 "hellowiki/common/utils"
	"log"
	"strconv"
)

type Role struct {
	gorm.Model
	RoleName string `gorm:"type:varchar(32)" json:"roleName"`
}

func GetAllRoles(pageNum int, pageSize int) (roles []Role, total int64) {
	dbBase := utils2.OpenDB()
	err := dbBase.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&roles).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return roles, 0
	}
	return roles, total
}

func GetRolesForUserById(userId int) (res []string, code int) {
	enforcer := utils2.GetEnforcer()
	if enforcer == nil {
		log.Printf("创建鉴权器失败")
		return res, result.ERROR
	}
	userIdStr := strconv.Itoa(userId)
	res, err := enforcer.GetRolesForUser(userIdStr)
	if err != nil {
		log.Printf("获取用户角色失败:{%v}", err)
		return res, result.ERROR
	}
	return res, result.SUCCSE
}

// 插入用户数据
func InsertRoleInfo(data Role) (code int) {
	dbBase := utils2.OpenDB()
	err := dbBase.Create(&data).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

func HasRoleByName(roleName string) (code int) {
	var role Role
	dbBase := utils2.OpenDB()
	dbBase.Take(&role, "role_name=?", roleName)
	if role.ID > 0 {
		//角色已存在
		return result.ERROR_ROLE_EXIST
	}
	//用户不存在
	return result.SUCCSE
}
