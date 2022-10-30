package service

import (
	"golang.org/x/crypto/bcrypt"
	"hellowiki/model"
	"log"
)

func CreateUser(userInfo *model.RegUser) (code int) {
	codeHas := model.HasUser(userInfo.UserName)
	if codeHas != 200 {
		return codeHas
	}
	//密码加密
	userInfo.PassWord = pswCrypt(userInfo.PassWord)
	codeInsert := model.Insert(userInfo)
	if codeInsert != 200 {
		return codeInsert
	}
	return codeInsert
}

func GetAllRegUserInfo(pageSize int, pageNum int) []model.RegUser {
	return model.FindAll(pageSize, pageNum)
}

func DeleteUser(userId int) int {
	return model.DeleteById(userId)
}

// 密码加密
func pswCrypt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
