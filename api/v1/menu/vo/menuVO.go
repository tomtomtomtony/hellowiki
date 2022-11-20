package vo

type MenuVO struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	ParentId   uint   `json:"parentId"`
	ParentName string `json:"parentName"`
	Type       string `json:"type"`
}

type TopMenu struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	ParentMenuId uint   `json:"parentMenuId"`
	ParentName   string `json:"parentName"`
}
