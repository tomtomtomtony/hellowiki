package result

/*
*
返回码
*/
const (
	SUCCSE = 200
	ERROR  = 500
	//code=1000....用户模块错误
	ERROR_USERNAME_USED  = 1001
	ERROR_PASSWORD_WRONG = 1002
	ERROR_USER_NOT_FOUND = 1003
)

var codeMsg = map[int]string{
	SUCCSE:               "OK",
	ERROR:                "FAIL",
	ERROR_USERNAME_USED:  "用户已存在",
	ERROR_USER_NOT_FOUND: "用户不存在",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
