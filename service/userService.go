package service

import "hellowiki/model"

func CreateUser(userInfo *model.RegUser) (code int) {
	codeHas := model.HasUser(userInfo.UserName)
	if codeHas != 200 {
		return codeHas
	}
	codeInsert := model.Insert(userInfo)
	if codeInsert != 200 {
		return codeInsert
	}
	return codeInsert
}
