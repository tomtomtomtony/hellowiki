package vo

type RegUserVO struct {
	UserName   string `json:"userName" validate:"required"`
	PassWord   string `json:"passWord" validate:"required"`
	ConfirmPwd string `json:"confirmPwd" validate:"required"`
}
