package vo

import "container/list"

type ConditionVO struct {
	ArticleId      uint      `json:"articleId" `
	CategoryName   string    `json:"categoryName"`
	ArticleTitle   string    `json:"articleTitle"`
	ArticleContent string    `json:"articleContent"`
	PageSize       int       `json:"pageSize"`
	PageNum        int       `json:"pageNum"`
	ParentLevel    uint      `json:"parentLevel"`
	Path           string    `json:"path"`
	ParentPath     string    `json:"parentPath"`
	Author         string    `json:"author"`
	Keywords       list.List `json:"keywords"`
}

type ResultVo struct {
	ArticleId uint   `json:"articleId" `
	Content   string `json:"content"`
	Title     string `json:"title"`
	Author    string `json:"author"`
}

type ResultArticle struct {
	Title          string `json:"title"`
	Id             uint   `json:"id"`
	CreateAt       int64  `json:"createAt"`
	UpdateAt       int64  `json:"updateAt"`
	CategoryMenuId uint   `json:"categoryMenuId"`
	CategoryName   string `json:"categoryName"`
}
