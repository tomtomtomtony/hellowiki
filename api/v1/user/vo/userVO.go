package vo

type RegUserVO struct {
	UserName   string `json:"userName" validate:"required"`
	PassWord   string `json:"passWord" validate:"required"`
	ConfirmPwd string `json:"confirmPwd" validate:"required"`
	IsEnable   bool   `json:"isEnable"`
	Roles      string `json:"roles"`
}

type LoginUserVO struct {
	Code     int    `json:"code"`
	UserName string `json:"userName" validate:"required"`
	PassWord string `json:"passWord" validate:"required"`
	Token    string `json:"token"`
}

type UserResult struct {
	UserName string `json:"userName"`
	Id       uint   `json:"id"`
	CreateAt int64  `json:"createAt"`
	UpdateAt int64  `json:"updateAt"`
	Roles    string `json:"roles"`
}

type UserList struct {
	Total    int64        `json:"total"`
	UserList []UserResult `json:"userList"`
}
