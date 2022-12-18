package vo

type RoleConditionVO struct {
	RoleName string   `json:"roleName"`
	Roles    []string `json:"roles"`
	PageSize int      `json:"pageSize"`
	PageNum  int      `json:"pageNum"`
}

type RoleResult struct {
	RoleName string `json:"roleName"`
	Id       uint   `json:"id"`
	CreateAt int64  `json:"createAt"`
	UpdateAt int64  `json:"updateAt"`
}

type RoleList struct {
	Total    int64        `json:"total"`
	RoleList []RoleResult `json:"roleList"`
}
