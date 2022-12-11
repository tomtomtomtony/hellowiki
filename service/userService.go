package service

import (
	"golang.org/x/crypto/bcrypt"
	"hellowiki/api/result"
	"hellowiki/api/v1/user/vo"
	"hellowiki/common/utils"
	"hellowiki/model"
	"log"
	"reflect"
)

func UserLogin(condition vo.LoginUserVO) (userName string, token string, code int) {
	userInfo := model.FindByName(condition.UserName)
	if reflect.DeepEqual(userInfo, model.RegUser{}) {
		log.Println("用户不存在")
		return "", "", result.ERROR_USER_NOT_FOUND
	}
	inputPwd := []byte(condition.PassWord)
	dbFindPwd := []byte(userInfo.PassWord)
	err := bcrypt.CompareHashAndPassword(dbFindPwd, inputPwd)
	if err != nil {
		log.Println("密码不正确")
		return "", "", result.ERROR_PASSWORD_WRONG
	}
	token = utils.NewJWT().IssueToken(userInfo.UserName)
	return userInfo.UserName, token, result.SUCCSE
}

func CreateUser(condition vo.RegUserVO) (code int) {
	var userInfo = vo2Do(condition)
	codeHas := model.HasUserByName(userInfo.UserName)
	if codeHas != result.SUCCSE {
		return codeHas
	}
	//密码加密
	userInfo.PassWord = pswCrypt(userInfo.PassWord)
	//默认启用账号
	userInfo.IsEnable = true
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

func SetUserName(id uint, condition vo.RegUserVO) int {
	if model.HasUserById(id) == result.SUCCSE {
		return result.ERROR_USER_NOT_FOUND
	}
	return model.UpdateUserById(id, "user_name", condition.UserName)
}

func SetRoles(id uint, condition vo.RegUserVO) int {
	if model.HasUserById(id) == result.SUCCSE {
		return result.ERROR_USER_NOT_FOUND
	}
	return model.UpdateUserById(id, "roles", condition.Roles)
}

func regUserVo2Do(vo vo.RegUserVO) model.RegUser {
	var res model.RegUser
	res.UserName = vo.UserName
	res.PassWord = vo.PassWord
	res.IsEnable = vo.IsEnable
	res.Roles = vo.Roles
	return res
}
