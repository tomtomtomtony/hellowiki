package vo

type CategoryVO struct {
	Name       string `json:"name"`
	ParentName string `json:"parentName"`
	ParentPath string `json:"parentPath"`
	Type       string `json:"type"`
	PageSize   int    `json:"pageSize"`
	PageNum    int    `json:"pageNum"`
	Path       string `json:"path"`
}
