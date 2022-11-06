package vo

type ConditionVO struct {
	CategoryId uint   `json:"categoryId"`
	Name       string `json:"name"`
	EngName    string `json:"engName"`
	ParentId   uint   `json:"parentId"`
	ParentName string `json:"parentName"`
	PageSize   int    `json:"pageSize"`
	PageNum    int    `json:"pageNum"`
}
