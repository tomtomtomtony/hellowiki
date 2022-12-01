package model

import (
	"gorm.io/gorm"
	"hellowiki/api/result"
	utils2 "hellowiki/common/utils"
)

type RegUser struct {
	gorm.Model
	UserName string `gorm:"type:varchar(20);not null " json:"userName"`
	PassWord string `gorm:"type:varchar(30);not null " json:"passWord"`
	IsEnable bool   `gorm:"type:boolean;not null" json:"isEnable"`
}

func HasUserById(id uint) (code int) {
	var regUser RegUser
	dbBase := utils2.OpenDB()
	dbBase.Take(&regUser, "id=?", id)
	if regUser.ID > 0 {
		//用户已存在
		return result.ERROR_USERNAME_USED
	}
	//用户不存在
	return result.SUCCSE
}

func HasUserByName(userName string) (code int) {
	var regUser RegUser
	dbBase := utils2.OpenDB()
	dbBase.Take(&regUser, "user_name=?", userName)
	if regUser.ID > 0 {
		//用户已存在
		return result.ERROR_USERNAME_USED
	}
	//用户不存在
	return result.SUCCSE
}

// 插入用户数据
func CreateUser(data RegUser) (code int) {
	dbBase := utils2.OpenDB()
	err := dbBase.Create(&data).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 查询用户列表
func FindAllUser(pageSize int, pageNum int) []RegUser {
	var users []RegUser
	dbBase := utils2.OpenDB()
	err := dbBase.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []RegUser{}
	}
	return users
}

// 根据id，软删除用户信息
func DeleteUserById(id int) int {
	dbBase := utils2.OpenDB()
	err := dbBase.Delete(&RegUser{}, "id=?", id).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 根据id，更新用户信息
func UpdateUserById(id uint, regUser RegUser) int {
	dbBase := utils2.OpenDB()
	err := dbBase.Model(&regUser).Where("id=?", id).Updates(regUser).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 按username查询
func FindByName(username string) RegUser {
	var regUser RegUser
	dbBase := utils2.OpenDB()
	err := dbBase.Take(&regUser, "user_name=?", username).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return RegUser{}
	}
	return regUser
}
