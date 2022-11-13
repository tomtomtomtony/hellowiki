package vo

type MenuVO struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	ParentId   uint   `json:"parentId"`
	ParentName string `json:"parentName"`
	Type       uint8  `json:"type"`
}

type TopMenu struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
