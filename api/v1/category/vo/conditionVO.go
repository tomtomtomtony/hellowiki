package vo

type ConditionVO struct {
	MenuId       uint   `json:"menuId"`
	Name         string `json:"name"`
	ParentMenuId uint   `json:"parentMenuId"`
	ParentName   string `json:"parentName"`
	PageSize     int    `json:"pageSize"`
	PageNum      int    `json:"pageNum"`
	ParentLevel  uint   `json:"parentLevel"`
	Path         string `json:"path"`
}
