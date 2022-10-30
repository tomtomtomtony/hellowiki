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

func HasUser(userName string) (code int) {
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
