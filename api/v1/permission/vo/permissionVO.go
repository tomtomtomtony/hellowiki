package vo

type PermissionVO struct {
	PermissionName string `json:"permissionName"`
}

type PermissionResult struct {
	PermissionName string `json:"permissionName"`
	Id             uint   `json:"id"`
	CreateAt       int64  `json:"createAt"`
	UpdateAt       int64  `json:"updateAt"`
}
type PermissionList struct {
	Total          int64              `json:"total"`
	PermissionList []PermissionResult `json:"permissionList"`
}
