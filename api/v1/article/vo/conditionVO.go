package vo

type ConditionVO struct {
	ArticleId       uint   `json:"articleId" `
	CategoryId      uint   `json:"categoryId"`
	CategoryName    string `json:"categoryName"`
	CategoryEngName string `json:"categoryEngName"`
	ArticleTitle    string `json:"articleTitle"`
	ArticleContent  string `json:"articleContent"`
	PageSize        int    `json:"pageSize"`
	PageNum         int    `json:"pageNum"`
}
