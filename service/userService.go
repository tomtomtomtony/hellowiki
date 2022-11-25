package service

import (
	"golang.org/x/crypto/bcrypt"
	"hellowiki/api/result"
	"hellowiki/api/v1/user/vo"
	"hellowiki/model"
	"log"
)

func CreateUser(condition vo.RegUserVO) (code int) {
	var userInfo = vo2Do(condition)
	codeHas := model.HasUserByName(userInfo.UserName)
	if codeHas != result.SUCCSE {
		return codeHas
	}
	//密码加密
	userInfo.PassWord = pswCrypt(userInfo.PassWord)
	codeInsert := model.CreateUser(userInfo)
	if codeInsert != result.SUCCSE {
		return codeInsert
	}
	return codeInsert
}
func vo2Do(userVo vo.RegUserVO) model.RegUser {
	var do model.RegUser
	do.UserName = userVo.UserName
	do.PassWord = userVo.PassWord
	return do
}
func GetAllRegUserInfo(pageSize int, pageNum int) []model.RegUser {
	return model.FindAllUser(pageSize, pageNum)
}

func DeleteUser(userId int) int {
	return model.DeleteUserById(userId)
}

// 密码加密
func pswCrypt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func SetUser(id uint, condition model.RegUser) int {
	if model.HasUserById(id) == result.SUCCSE {
		return result.ERROR_USER_NOT_FOUND
	}
	return model.UpdateUserById(id, condition)
}
