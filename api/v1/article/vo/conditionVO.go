package vo

type ConditionVO struct {
	ArticleId      uint   `json:"articleId" `
	CategoryMenuId uint   `json:"categoryMenuId"`
	CategoryName   string `json:"categoryName"`
	ArticleTitle   string `json:"articleTitle"`
	ArticleContent string `json:"articleContent"`
	PageSize       int    `json:"pageSize"`
	PageNum        int    `json:"pageNum"`
	ParentLevel    uint   `json:"parentLevel"`
	Path           string `json:"path"`
}
