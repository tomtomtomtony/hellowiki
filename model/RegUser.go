package model

import (
	"gorm.io/gorm"
	"hellowiki/common/result"
)

type RegUser struct {
	gorm.Model
	UserName string `gorm:"type:varchar(20);not null " json:"userName"`
	PassWord string `gorm:"type:varchar(30);not null " json:"passWord"`
}

func HasUserById(id uint) (code int) {
	var regUser RegUser
	Db.Take(&regUser, "id=?", id)
	if regUser.ID > 0 {
		//用户已存在
		return result.ERROR_USERNAME_USED
	}
	//用户不存在
	return result.SUCCSE
}

func HasUserByName(userName string) (code int) {
	var regUser RegUser
	Db.Take(&regUser, "user_name=?", userName)
	if regUser.ID > 0 {
		//用户已存在
		return result.ERROR_USERNAME_USED
	}
	//用户不存在
	return result.SUCCSE
}

// 插入用户数据
func Insert(data *RegUser) (code int) {
	err := Db.Create(&data).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 查询用户列表
func FindAll(pageSize int, pageNum int) []RegUser {
	var users []RegUser
	err := Db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	return users
}

// 根据id，软删除用户信息
func DeleteById(id int) int {
	err := Db.Delete(&RegUser{}, "id=?", id).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 根据id，更新用户信息
func UpdateById(id uint, regUser RegUser) int {
	err := Db.Model(&regUser).Where("id=?", id).Updates(regUser).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 按username查询
//func FindByName(username string) RegUser {
//	var regUser RegUser
//	err := Db.Take(&regUser, "user_name=?", username).Error
//	if err != nil && err != gorm.ErrRecordNotFound {
//		return
//	}
//	return
//}
